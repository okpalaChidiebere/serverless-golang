package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-xray-sdk-go/xray"
	uuid "github.com/satori/go.uuid"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type Image struct {
	ImageId   string `json:"imageId"`
	GroupId   string `json:"groupId"`
	Title     string `json:"title"`
	Timestamp string `json:"timestamp"`
	ImageUrl  string `json:"imageUrl"`
}

type createImageResponse struct {
	Image     Image  `json:"newItem"`
	UploadUrl string `json:"uploadUrl"`
}

var (
	ddb *dynamodb.DynamoDB
	s3s *s3.S3
	gTb = aws.String(os.Getenv("GROUPS_TABLE"))
	iTb = aws.String(os.Getenv("IMAGES_TABLE"))
	bn  = os.Getenv("IMAGES_S3_BUCKET")
)

func init() {
	svc := session.Must(session.NewSession()) // Use aws sdk to connect to dynamoDB
	ddb = dynamodb.New(svc)                   // Create DynamoDB client
	s3s = s3.New(svc)                         // Create S3 service client

	xray.AWS(ddb.Client)
	xray.AWS(s3s.Client)
}

func createImageHandler(req Request) (Response, error) {
	var buf bytes.Buffer

	r, _ := json.MarshalIndent(req, "", " ")
	log.Printf("Caller Request: %s", r)

	// Parse groupId variable from request url
	gId := req.PathParameters["groupId"]

	c := make(chan bool)
	nIc := make(chan Image)

	go groupExists(gId, c)

	validGroupId := <-c

	if !validGroupId {

		body, _ := json.Marshal(map[string]interface{}{
			"error": "Group does not exist",
		})
		json.HTMLEscape(&buf, body)

		return Response{
			StatusCode: 404,
			Body:       buf.String(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	/* We are sure that our group does exist at this point */
	imageId := uuid.Must(uuid.NewV4(), nil).String() //create a new id

	go createImage(gId, imageId, req, nIc)

	nIt := <-nIc

	url, _ := getUploadUrl(imageId)

	body, _ := json.Marshal(&createImageResponse{
		nIt,
		url,
	})
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

func groupExists(gId string, c chan bool) {
	// Build the query input parameters
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(gId),
			},
		},
		TableName: gTb,
	}

	// Make the DynamoDB Query API call
	//https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB.GetItem
	rslt, _ := ddb.GetItem(params)

	if rslt.Item != nil {
		c <- true // we know group exist
	} else {
		c <- false // if not item is returned, we know the group does not exist
	}
}

func createImage(groupId string, imageId string, event Request, c chan Image) {
	// Initialize group
	newItem := &Image{
		ImageId:   imageId,
		Timestamp: time.Now().String(),
		GroupId:   groupId,
		ImageUrl:  "https://" + bn + ".s3.amazonaws.com/" + imageId,
	}

	// Parse request body
	json.Unmarshal([]byte(event.Body), newItem)

	e, _ := json.MarshalIndent(newItem, "", " ")
	log.Printf("Storing new item: %s", e)

	// Write the new item to DynamoDB database
	item, _ := dynamodbattribute.MarshalMap(newItem)
	p := &dynamodb.PutItemInput{
		Item:      item,
		TableName: iTb,
	}

	if _, err := ddb.PutItem(p); err != nil {
		log.Fatalf("Failed to create new item: Error message was %s", err.Error())
	}

	c <- *newItem
}

func getUploadUrl(imageId string) (string, error) {
	//More on PutObject here https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.PutObjectRequest
	req, _ := s3s.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bn),
		Key:    aws.String(imageId),
	})

	urlStr, err := req.Presign(5 * time.Minute)

	return urlStr, err
}

func main() {
	lambda.Start(createImageHandler)
}
