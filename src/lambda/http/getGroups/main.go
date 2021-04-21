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
)

type Response events.APIGatewayProxyResponse

type Group struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetGroupsResponse struct {
	Groups []Group `json:"items"`
}

var ddb *dynamodb.DynamoDB
var (
	tableName = aws.String(os.Getenv("GROUPS_TABLE"))
)

func init() {
	session := session.Must(session.NewSession()) // Use aws sdk to connect to dynamoDB
	ddb = dynamodb.New(session)                   // Create DynamoDB client
}

func GetGroupsHandler(ctx context.Context) (Response, error) {
	log.Println("GetGroups")
	var buf bytes.Buffer

	// Read from DynamoDB
	input := &dynamodb.ScanInput{
		TableName: tableName,
	}
	result, _ := ddb.Scan(input)

	// Construct todos from response
	var groups []Group
	for _, i := range result.Items {
		group := Group{}
		if err := dynamodbattribute.UnmarshalMap(i, &group); err != nil {
			log.Println("Failed to unmarshal")
			log.Println(err)
		}
		groups = append(groups, group)
	}

	// Success HTTP response
	body, err := json.Marshal(&GetGroupsResponse{
		groups,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode: 200,
		Body:       buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(GetGroupsHandler)
}
