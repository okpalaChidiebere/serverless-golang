package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/udacity/serverless-golang/src/businessLogic/groups"
	"github.com/udacity/serverless-golang/src/dataLayer/groupsAccess"
	"github.com/udacity/serverless-golang/src/models"
	"github.com/udacity/serverless-golang/src/requests"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type CreateGroupResponse struct {
	Group models.Group `json:"newItem"`
}

func createGroupHandler(ctx context.Context, req Request) (Response, error) {
	e, _ := json.MarshalIndent(req, "", " ")
	log.Printf("Processing Event: %s", e)

	var buf bytes.Buffer

	// Initialize CreateGroupRequest
	group := &requests.CreateGroupRequest{}

	groupsRepo := groupsAccess.NewDynamoDbRepo()
	ga := groups.NewGroupAccess(groupsRepo)

	// Parse request body
	json.Unmarshal([]byte(req.Body), group)

	newItem, err := ga.CreateGroup(group)
	if err != nil {

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
			newItem,
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
