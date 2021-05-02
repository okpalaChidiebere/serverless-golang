package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type Image struct {
	ImageId   string `json:"imageId"`
	GroupId   string `json:"groupId"`
	Title     string `json:"title"`
	Timestamp string `json:"timestamp"`
}

var (
	ddb          *dynamodb.DynamoDB
	imageIdTable = aws.String(os.Getenv("IMAGE_ID_INDEX"))
	imageTable   = aws.String(os.Getenv("IMAGES_TABLE"))
)

func init() {
	svc := session.Must(session.NewSession()) // Use aws sdk to connect to dynamoDB
	ddb = dynamodb.New(svc)                   // Create DynamoDB client
	xray.AWS(ddb.Client)
}

func getImageHandler(req Request) (Response, error) {
	var buf bytes.Buffer

	r, _ := json.MarshalIndent(req, "", " ")
	log.Printf("Caller Request: %s", r)

	// Parse groupId variable from request url
	mId := req.PathParameters["imageId"]

	c := make(chan []map[string]*dynamodb.AttributeValue)

	go func() {
		p := &dynamodb.QueryInput{
			TableName:              imageTable,
			IndexName:              imageIdTable,
			KeyConditionExpression: aws.String("imageId = :imageId"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":imageId": {
					S: aws.String(mId),
				},
			},
		}

		// Make the DynamoDB Query API call
		rslt, _ := ddb.Query(p)
		c <- rslt.Items
	}() //imageId string

	results := <-c

	if len(results) > 0 {
		item := Image{}

		dynamodbattribute.UnmarshalMap(results[0], &item)
		body, _ := json.Marshal(item)
		json.HTMLEscape(&buf, body)

		return Response{
			StatusCode: 200,
			Body:       buf.String(),
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return Response{
		StatusCode: 404,
		Body:       "",
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	lambda.Start(getImageHandler)
}
