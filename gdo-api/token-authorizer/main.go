package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Help function to generate an IAM policy
func generatePolicy(principalId, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	// Optional output with custom properties of the String, Number or Boolean type.
	authResponse.Context = map[string]interface{}{
		"stringKey":  "stringval",
		"numberKey":  123,
		"booleanKey": true,
	}
	return authResponse
}

func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := event.AuthorizationToken
	parsed, err := jwt.ParseSigned(token) // decode token
	if err != nil {
		// for now authorize
		fmt.Println(err)
		return generatePolicy("user", "Allow", event.MethodArn), nil
		//return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
	}
	c := make(map[string]interface{})
	pubKey := os.Getenv("PUBKEY")
	err = parsed.Claims(&pubKey, &c)
	if err != nil {
		// for now authorize
		fmt.Println(err)
		return generatePolicy("user", "Allow", event.MethodArn), nil
		//return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
	}
	//test claims
	if false {
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}
	// check authority
	if false {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
	}
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
	}
	// for now authorize
	fmt.Printf("auth - no errors iss: %s, sub: %s\n", c["iss"], c["sub"])
	return generatePolicy("user", "Allow", event.MethodArn), nil

}

func main() {
	lambda.Start(handleRequest)
}
