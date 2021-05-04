package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/udacity/serverless-golang/src/businessLogic/groups"
	"github.com/udacity/serverless-golang/src/dataLayer/groupsAccess"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

func GetGroupsHandler(ctx context.Context, req Request) (Response, error) {
	log.Println("GetGroups")
	var buf bytes.Buffer

	queryParams := req.QueryStringParameters
	if len(queryParams) == 0 {
		log.Println("Undefined queryString parameters")
		//the client gets 502(Bad Gateway) error returned but the error returned will be logged to cloudwatch
		return Response{StatusCode: 400}, errors.New("bad request")
	}

	var nextKey string // Next key to continue scan operation if necessary
	nk, ok := queryParams["nextKey"]
	if !ok {
		nextKey = ""
	} else {
		nextKey = nk
	}

	var limit int64 // Maximum number of elements to return
	if ql, ok := queryParams["limit"]; !ok {
		limit = 20 //default limit is 20 if not lmit is given
	} else {
		//convert string to int64
		fmt.Sscan(ql, &limit)
	}

	if limit <= 0 {
		log.Println("Limit parameter should be positive")
		return Response{StatusCode: 400}, errors.New("bad request")
	}

	groupsRepo := groupsAccess.NewDynamoDbRepo()
	ga := groups.NewGroupAccess(groupsRepo)

	groups, nk := ga.GetAllGroups(limit, nextKey)

	// Success HTTP response
	body, err := json.Marshal(map[string]interface{}{
		"items":   groups,
		"nextKey": url.QueryEscape(nk),
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
