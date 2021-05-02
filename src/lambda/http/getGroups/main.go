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
)

type Response events.APIGatewayProxyResponse

type GetGroupsResponse struct {
	Groups []models.Group `json:"items"`
}

func GetGroupsHandler(ctx context.Context) (Response, error) {
	log.Println("GetGroups")
	var buf bytes.Buffer

	groupsRepo := groupsAccess.NewDynamoDbRepo()
	ga := groups.NewGroupAccess(groupsRepo)

	groups := ga.GetAllGroups()

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
