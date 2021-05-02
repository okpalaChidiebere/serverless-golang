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

type getImagesResponse struct {
	Images []Image `json:"items"`
}

var (
	ddb *dynamodb.DynamoDB
	gTb = aws.String(os.Getenv("GROUPS_TABLE"))
	iTb = aws.String(os.Getenv("IMAGES_TABLE"))
)

func init() {
	svc := session.Must(session.NewSession()) // Use aws sdk to connect to dynamoDB
	ddb = dynamodb.New(svc)                   // Create DynamoDB client
	xray.AWS(ddb.Client)
}

func getImagesHandler(req Request) (Response, error) {
	var buf bytes.Buffer

	r, _ := json.MarshalIndent(req, "", " ")
	log.Printf("Caller Request: %s", r)

	// Parse groupId variable from request url
	gId := req.PathParameters["groupId"]

	c := make(chan bool)
	ic := make(chan []Image)

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
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	/* We are sure that our group does exist at this point */
	go getImagesPerGroup(gId, ic)

	imgs := <-ic

	// Success HTTP response
	body, _ := json.Marshal(&getImagesResponse{
		Images: imgs,
	})
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode: 200,
		Body:       buf.String(),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}

	return resp, nil
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

func getImagesPerGroup(gId string, c chan []Image) {

	p := &dynamodb.QueryInput{
		TableName:              iTb,
		KeyConditionExpression: aws.String("groupId = :groupId"), //we specify that we want all the images that has the the groupId we want
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":groupId": {
				S: aws.String(gId),
			},
			//more on how to use other atttibute types here https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue#Unmarshal
		},
		ScanIndexForward: aws.Bool(false), //it reverses the order of the list. The latest images will be first
	}

	// Make the DynamoDB Query API call
	rslt, _ := ddb.Query(p)

	// Construct todos from response
	var imgs []Image
	for _, i := range rslt.Items {
		img := Image{}
		if err := dynamodbattribute.UnmarshalMap(i, &img); err != nil {
			log.Println("Failed to unmarshal")
			log.Println(err)
		}
		imgs = append(imgs, img)
	}

	c <- imgs
}

func main() {
	lambda.Start(getImagesHandler)
}
