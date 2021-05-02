package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayWebsocketProxyRequest

var ddb *dynamodb.DynamoDB
var (
	ct = aws.String(os.Getenv("CONNECTIONS_TABLE"))
)

func init() {
	session := session.Must(session.NewSession()) // Use aws sdk to connect to dynamoDB
	ddb = dynamodb.New(session)                   // Create DynamoDB client

	xray.AWS(ddb.Client)
}

func disconnectHandler(req Request) (Response, error) {

	r, _ := json.MarshalIndent(req, "", " ")
	log.Printf("Websocket dicconnect: %s", r)

	// Parse connectionID from websocketrequest url
	cId := req.RequestContext.ConnectionID

	//We want to delete a connection by its id
	p := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(cId),
			},
		},
		TableName: ct,
	}

	if _, err := ddb.DeleteItem(p); err != nil {
		log.Fatalf("Failed to delete user socket map: Error message was %s", err.Error())
	}

	return Response{
		StatusCode: 200,
		Body:       "",
	}, nil
}

func main() {
	lambda.Start(disconnectHandler)
}
