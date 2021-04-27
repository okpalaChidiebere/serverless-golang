package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type s3Event events.S3Event

type Payload struct {
	ImageId string `json:"imageId"`
}

type UserConn struct {
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

var ddb *dynamodb.DynamoDB
var apiGateway *apigatewaymanagementapi.ApiGatewayManagementApi
var (
	ct    = os.Getenv("CONNECTIONS_TABLE")
	stage = os.Getenv("STAGE")
	apiId = os.Getenv("API_ID")
)

func init() {
	session := session.Must(session.NewSession()) // Use aws sdk to connect to dynamoDB
	ddb = dynamodb.New(session)                   // Create DynamoDB client
	apiGateway = apigatewaymanagementapi.New(session, aws.NewConfig().WithEndpoint(apiId+".execute-api.ca-central-1.amazonaws.com/"+stage))
}

func sendNotificationsHandler(event s3Event) {
	for _, record := range event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)

		c := make(chan string)

		// Read the list of connected users(IDs) from DynamoDB
		params := &dynamodb.ScanInput{
			TableName: aws.String(ct),
		}
		result, _ := ddb.Scan(params)

		/*
		   For simplicity we will just send an ID of an image that was uploaded as payload.

		   We could do more complicated logic and send more data in payload like
		   fetch the image information based on the imageId. We have an apiGatway that we can invoke
		*/
		p := Payload{s3.Object.Key}

		conns := make([]UserConn, *result.Count)
		dynamodbattribute.UnmarshalListOfMaps(result.Items, &conns)

		for _, i := range conns {
			go sendMessageToClient(i.Id, p, c)
		}

		//make sure all the notification messages are sent
		for i := 0; i < len(result.Items); i++ {
			fmt.Printf("Done sending message to client connection id: %s", <-c)
		}
	}
}

func sendMessageToClient(connId string, pload Payload, c chan string) {
	fmt.Println("Sending message to a connection", connId)

	body, _ := json.Marshal(pload)

	connParams := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connId),
		Data:         body,
	}
	_, err := apiGateway.PostToConnection(connParams)
	if err != nil {
		//check tthe error type returned. More in the link:
		//https://docs.aws.amazon.com/sdk-for-go/api/service/apigatewaymanagementapi/#ApiGatewayManagementApi.PostToConnection
		switch terr := err.(type) {
		case *apigatewaymanagementapi.GoneException:
			if terr.StatusCode() == 410 { //we still have connectionId in our db but that connection was closed
				fmt.Println("Stale connection")

				// Delete this connection from the db table
				_, err := ddb.DeleteItem(&dynamodb.DeleteItemInput{
					Key: map[string]*dynamodb.AttributeValue{
						"id": {
							S: aws.String(connId),
						},
					},
					TableName: aws.String(ct),
				})

				if err != nil {
					log.Fatalf(err.Error())
				}
			}
		default:
			fmt.Printf("unhandled error of type %T: %s", err, err)

		}

		/*
			//Another way to handle error type instead of using switch is below but i think the one i wrote is better
			//for debugging purposes
			//More ways of error handling in go https://blog.golang.org/go1.13-errors
			//https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/handling-errors.html

			if tErr, ok := err.(*apigatewaymanagementapi.GoneException); ok {
				if tErr.StatusCode() == 410 {

				}
			}else{
				//process error generally
			}
		*/
	}

	c <- connId
}

func main() {
	lambda.Start(sendNotificationsHandler)
}
