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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type DynamoDBStreamEvent events.DynamoDBEvent

type Body struct {
	ImageId   string `json:"imageId"`
	GroupId   string `json:"groupId"`
	Title     string `json:"title"`
	Timestamp string `json:"timestamp"`
	ImageUrl  string `json:"imageUrl"`
}

type AWSUser struct {
	// in the aws secretManager resource, we can have multiple fields which each containing a secret value. We will have just two fields
	KeyID     string `json:"AWS_ACCESS_KEY_ID"`
	SecretKey string `json:"AWS_SECRET_ACCESS_KEY"`
}

var (
	signer   *v4.Signer
	url      string
	esHost   = fmt.Sprintf("https://%s", os.Getenv("ES_ENDPOINT"))
	index    = "images-index"
	t        = "images" //the type, as required by ES as the path to your indicies
	service  = "es"
	region   = "ca-central-1"
	sm       *secretsmanager.SecretsManager
	secretId = os.Getenv("AWS_APP_USER_SECRET_ID")
	/*
		The point of this variable is to cache the value of our secret and dont sent requst to
		secretManager over and over again
		AWS Lambda may keep our function instance for sometime, then we will reuse this cached secret
		and we wil save some money by not calling secretManager
	*/
	cachedSecretObj *AWSUser
)

func init() {

	c := make(chan *AWSUser)

	s := session.Must(session.NewSession()) // Use aws sdk
	sm = secretsmanager.New(s)              // Create SecretsManager client

	go getSecret(c)

	secretObj := <-c

	// Get credentials from environment variables and create the AWS Signature Version 4 signer
	credentials := credentials.NewStaticCredentials(secretObj.KeyID, secretObj.SecretKey, "")
	signer = v4.NewSigner(credentials)

	url = esHost + "/" + index + "/" + t + "/"
}

func getSecret(c chan *AWSUser) {

	//If there is a cached secret already then we just return it
	if cachedSecretObj != nil {
		c <- cachedSecretObj
		return
	}

	data, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		c <- &AWSUser{}
		return
	}

	// Initialize userScret
	u := &AWSUser{}

	// we parse the *secretstring because it is a JSON object. eg "{\n  \"AWS_ACCESS_KEY_ID\":\"someVale\"\n}\n" needs to be parsed to get the JSON object
	//NOTE: we dereference the pointer to a string, then get the string so that we can cast it to []byte
	json.Unmarshal([]byte(*data.SecretString), u)

	cachedSecretObj = u
	c <- cachedSecretObj
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
