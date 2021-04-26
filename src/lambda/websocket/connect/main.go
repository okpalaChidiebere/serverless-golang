package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayWebsocketProxyRequest

type UserConn struct {
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

var ddb *dynamodb.DynamoDB
var (
	ct = os.Getenv("CONNECTIONS_TABLE")
)

func init() {
	session := session.Must(session.NewSession()) // Use aws sdk to connect to dynamoDB
	ddb = dynamodb.New(session)                   // Create DynamoDB client
}

func connectHandler(req Request) (Response, error) {

	r, _ := json.MarshalIndent(req, "", " ")
	log.Printf("Websocket connect: %s", r)

	// Parse connectionID from websocketrequest url
	cId := req.RequestContext.ConnectionID
	timestamp := time.Now().String() //or req.RequestContext.RequestTime

	conn := UserConn{
		cId,
		timestamp,
	}

	item, err := dynamodbattribute.MarshalMap(conn)
	if err != nil {
		log.Fatalf("Unable to marshal user socket map: Error message was %s", err.Error())
	}

	p := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(ct),
	}

	if _, err := ddb.PutItem(p); err != nil {
		log.Fatalf("Failed to create new item: Error message was %s", err.Error())
	}

	return Response{
		StatusCode: 200,
		Body:       "",
	}, nil
}

func main() {
	lambda.Start(connectHandler)
}
