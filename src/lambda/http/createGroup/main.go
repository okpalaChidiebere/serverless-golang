package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type Group struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateGroupResponse struct {
	Group *Group `json:"newItem"`
}

var ddb *dynamodb.DynamoDB
var (
	tableName = aws.String(os.Getenv("GROUPS_TABLE"))
)

func init() {
	session := session.Must(session.NewSession()) // Use aws sdk to connect to dynamoDB
	ddb = dynamodb.New(session)                   // Create DynamoDB client
}

func createGroupHandler(ctx context.Context, req Request) (Response, error) {
	//fmt.Println("Processing Event: ", req)
	log.Println("Processing Event: ", req)

	var buf bytes.Buffer
	id := uuid.Must(uuid.NewV4(), nil).String() //create a new id

	// Initialize group
	group := &Group{
		Id: id,
	}

	// Parse request body
	json.Unmarshal([]byte(req.Body), group)

	// Write the new item to DynamoDB database
	item, _ := dynamodbattribute.MarshalMap(group)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: tableName,
	}

	if _, err := ddb.PutItem(input); err != nil {

		log.Fatalf("Failed to create new item: Error message was %s", err.Error())
		// Error HTTP response
		resp := Response{
			StatusCode: 500,
			Body:       err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		return resp, nil
	} else {
		//body, _ := json.Marshal(group)
		body, _ := json.Marshal(&CreateGroupResponse{
			group,
		})
		json.HTMLEscape(&buf, body)

		resp := Response{
			StatusCode: 201,
			Body:       buf.String(),
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}

		return resp, nil
	}
}

func main() {
	lambda.Start(createGroupHandler)
}
