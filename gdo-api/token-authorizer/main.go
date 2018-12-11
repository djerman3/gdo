package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// global token decrypt key
var pubKey *jwk.Key

// Help function to generate an IAM policy
func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

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
	payload := event.AuthorizationToken
	token, err := jwt.ParseVerify(bytes.NewReader([]byte(payload)), jwa.RS256, pubKey)
	if err != nil {
		// invalid parseverify
		// for now authorize
		fmt.Println(err)
		return generatePolicy("user", "Allow", event.MethodArn), nil
		//return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
	}
	user := token.Subject()
	audience := token.Audience()
	expire := token.Expiration()
	// for now authorize
	fmt.Printf("Accepted token i:%s a:%s u:%s e:%s\n", token.Issuer(), audience, user, expire.Format(time.RFC3339))
	return generatePolicy("user", "Allow", event.MethodArn), nil
	//return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
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
	fmt.Printf("auth - no errors iss: %s, sub: %s\n", token.Issuer(), token.Subject())
	return generatePolicy("user", "Allow", event.MethodArn), nil

}

func main() {
	keyserver := os.Getenv("JWKS_URL")
	set, err := jwk.Fetch(keyserver)

	if err != nil {
		fmt.Println(err)
		return
	}
	keys := set.LookupKeyID(os.Getenv("KEY_ID"))

	if len(keys) == 0 {
		fmt.Println("no key by that id")
		keys = set.Keys
	}
	pubKey = &keys[0]
	lambda.Start(handleRequest)
}
