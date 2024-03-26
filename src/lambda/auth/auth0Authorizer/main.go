package main

import (
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Event events.APIGatewayCustomAuthorizerRequest
type Response events.APIGatewayCustomAuthorizerResponse

func main() {
	lambda.Start(auth0AuthorizerHandler)
}

func auth0AuthorizerHandler(event Event) (Response, error) {
	token := event.AuthorizationToken //We extract the token from the Auth Header

	if len(token) == 0 {
		log.Println("User was not authorized: No authentication header")
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}

	if s := strings.ToLower(token); !strings.HasPrefix(s, "bearer ") {
		log.Println("User was not authorized: Invalid authentication header")
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}

	//getting the value of the token from the header
	split := strings.Split(token, " ")
	bearerToken := split[1]

	/*
		Here we are checking if the token is not equal to a mock value we expect then we dont authorize the user

		Ideally this is where we will validate our real token from a third party service like Auth0, whether that token is a valid JWT token or not
		There is Symmetric and Asymmetric way

		The Symmetric way is you basically 'verify' the token with the secretKey. For Auth0, the secret is the 'Client Secret' from the dashboard
		token, err := jwt.Parse('ACCESS_TOKEN', func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHS256); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return 'CLIENT_SECRET', nil
		})

		The Asymmetric way is basically where you use Auth0 public cert or key to verify the token
		The key use to sign the accessToken is stored my Auth0; its private. We don't need to store this secret key ourself
		Eg: Auth0 public cert `curl https://AUTH0_DOMAIN.us.auth0.com/pem | openssl x509 -pubkey -noout`
		    Auth0 public key can be gotten by making a fetch request to https://AUTH0_DOMAIN.us.auth0.com/.well-known/jwks.json

			The code for verifying jwt signature with pem cert looks like
			var pubkey = `-----BEGIN PUBLIC KEY-----
			MIIBIjANB..........
			-----END PUBLIC KEY-----`
			mee, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubkey))
			if err != nil {
				log.Println("errorPublic:", err)
			}
			token, err := jwt.Parse('ACCESS_TOKEN', func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return mee, nil
			})

		    The code for verifying jwt signature with jwk can be seen here
			https://stackoverflow.com/questions/41077953/how-to-verify-jwt-signature-with-jwk-in-go

		These links will help you
		https://stackoverflow.com/questions/46735347/how-can-i-fetch-a-certificate-from-a-url
		https://stackoverflow.com/questions/66984610/problem-when-parsing-rs256-public-key-with-dgrijalva-jwt-go-golang-package
		https://auth0.com/docs/secure/tokens/access-tokens/get-management-api-access-tokens-for-testing  test auth0 token
		https://community.auth0.com/t/where-is-the-auth0-public-key-to-be-used-in-jwt-io-to-verify-the-signature-of-a-rs256-token/8455
		https://auth0.com/blog/authentication-in-golang/
		https://stackoverflow.com/questions/51834234/i-have-a-public-key-and-a-jwt-how-do-i-check-if-its-valid-in-go
		https://brunoscheufler.com/blog/2020-04-11-verifying-asymmetrically-signed-jwts-in-go
	*/
	if bearerToken != "123" {
		log.Println("User was not authorized: Invalid token")
		return Response{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
		//other ways to use new Error to have your function return an error in go https://www.geeksforgeeks.org/errors-new-function-in-golang-with-examples/
	}

	//At this point, there are no exceptions and the request has been authorized
	return generatePolicy("user", "Allow", event.MethodArn), nil
}

func generatePolicy(principalID, effect, resource string) Response {
	authResponse := Response{PrincipalID: principalID}

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
	//authResponse.Context = context //optional. I did not want to use it. More on this here https://github.com/aws/aws-lambda-go/blob/master/events/README_ApiGatewayCustomAuthorizer.md
	return authResponse
}
