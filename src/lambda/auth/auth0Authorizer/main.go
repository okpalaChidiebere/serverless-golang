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
		You basically 'verify' the token with the secretKey that Auth0 gives you

		These links will help you
		https://stackoverflow.com/questions/51834234/i-have-a-public-key-and-a-jwt-how-do-i-check-if-its-valid-in-go
		https://qvault.io/cryptography/how-to-build-jwts-in-go-golang/
		https://auth0.com/blog/authentication-in-golang/
		https://betterprogramming.pub/hands-on-with-jwt-in-golang-8c986d1bb4c0
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
