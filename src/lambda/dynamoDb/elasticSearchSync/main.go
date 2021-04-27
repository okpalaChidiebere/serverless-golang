package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type DynamoDBStreamEvent events.DynamoDBEvent

func elasticSearchSyncHandler(e DynamoDBStreamEvent) {

	r, _ := json.MarshalIndent(e, "", " ")
	log.Printf("Processing events batch from DynamoDB: %s", r)

	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)
	}

}

func main() {
	lambda.Start(elasticSearchSyncHandler)
}
