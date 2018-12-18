package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

const tokISSUER = "https://jerman3.auth0.com/"
const tokAUDIENCE = "https://gdo.heather.com"

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

	return authResponse
}

func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	// Find a key
	keyserver := os.Getenv("JWKS_URL")
	if len(keyserver) == 0 {
		// invalid url
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Error: Invalid keyserver var")
	}
	set, err := jwk.Fetch(keyserver)
	if err != nil {
		// invalid url
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Error: Invalid keyserver:%s", err)
	}
	// decode key
	keyId := os.Getenv("JWK_ID")
	keys := set.Keys
	if len(keyId) > 0 {
		keys = set.LookupKeyID(keyId)
		if len(keys) == 0 {
			return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Error: Invalid keyId:%s", keyId)
		}
	}
	key, err := set.Keys[0].Materialize()
	if err != nil {
		// invalid key
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Error: Invalid key:%s", err)
	}

	//parse
	payload := event.AuthorizationToken
	token, err := jwt.ParseVerify(bytes.NewReader([]byte(payload)), jwa.RS256, key)
	if err != nil {
		// invalid parseverify
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Error: Invalid token:%s", err)
	}
	if token.Issuer() != tokISSUER {
		// wrong claims
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}
	authResponse := generatePolicy("user", "Allow", event.MethodArn)
	// Optional output with custom properties of the String, Number or Boolean type.
	authResponse.Context = map[string]interface{}{
		"subject":  token.Subject(),
		"issuer":   token.Issuer(),
		"audience": token.Audience(),
		"issueAt":  token.IssuedAt(),
		"expires":  token.Expiration(),
	}
	return authResponse, nil
	//return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
	//test claims to return deny,err
	//if false {
	//	return generatePolicy("user", "Deny", event.MethodArn), nil
	//}
	// check authority
	//if false {
	//	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
	//	}
	//if err != nil {
	//	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
	//}

}

func main() {

	lambda.Start(handleRequest)
}
