package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

type DynamoDBStreamEvent events.DynamoDBEvent

type Body struct {
	ImageId   string `json:"imageId"`
	GroupId   string `json:"groupId"`
	Title     string `json:"title"`
	Timestamp string `json:"timestamp"`
	ImageUrl  string `json:"imageUrl"`
}

var (
	signer  *v4.Signer
	url     string
	esHost  = fmt.Sprintf("https://%s", os.Getenv("ES_ENDPOINT"))
	index   = "images-index"
	t       = "images" //the type, as required by ES as the path to your indicies
	service = "es"
	region  = "ca-central-1"
)

func init() {

	// Get credentials from environment variables and create the AWS Signature Version 4 signer
	credentials := credentials.NewStaticCredentials("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "") //TODO: Read the secret from secret manager
	signer = v4.NewSigner(credentials)

	url = esHost + "/" + index + "/" + t + "/"
}

func elasticSearchSyncHandler(e DynamoDBStreamEvent) error {

	r, _ := json.MarshalIndent(e, "", " ")
	log.Printf("Processing events batch from DynamoDB: %s", r)

	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

		/*
			If the eventName is not 'INSERT', we will skip this record.

			However, in a complete solution, we will remove record from ES if they are removed(REMOVE) from dynamodb
			or update ES record if they are updated(MODIFY) in dynamoDb. Eg like in the link
			https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_streams_Record.html
		*/
		if record.EventName != "INSERT" {
			continue
		}

		//Get the item that was added to dynaoDb
		newItem := record.Change.NewImage

		id := newItem["imageId"].String()

		/*
			We create the document that we want to store in ElasticSearch.
			We basically copied all fields from our dynamoDb item to our struct
		*/
		b := Body{
			ImageId:   newItem["imageId"].String(),
			Timestamp: newItem["timestamp"].String(),
			GroupId:   newItem["groupId"].String(),
			ImageUrl:  newItem["imageUrl"].String(),
			Title:     newItem["title"].String(),
		}

		// JSON document to be included as the request body
		out, _ := json.Marshal(b)
		body := strings.NewReader(string(out))

		// An HTTP client for sending the request
		client := &http.Client{}

		// Form the HTTP request. We are making a PUT requet as specified in ES doc
		//https://www.elastic.co/guide/en/elasticsearch/reference/6.8/docs-index_.html
		req, err := http.NewRequest(http.MethodPut, url+id, body)
		if err != nil {
			fmt.Print(err)
		}

		// You can probably infer Content-Type programmatically, but here, we just say that it's JSON
		req.Header.Add("Content-Type", "application/json")

		// Sign the request, send it, and print the response
		signer.Sign(req, body, service, region, time.Now())
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Print(resp.Status + "\n")

	}

	return nil
}

func main() {
	lambda.Start(elasticSearchSyncHandler)
}
