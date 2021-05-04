package groupsAccess

import (
	"encoding/json"
	"log"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/udacity/serverless-golang/src/models"
)

/*
This interface is a Port. It provides an interface that is independent from specific technologies

A Port can have multiple Adapters. This means any adapter that belongs to this Port, will implement
these two methods. In our app, we just have one Adapter(GroupDynamoDbRepository) that belongs to this Port(Repository)

Alternative to arrange your code; but i think mine is fine
https://github.com/yuraxdrumz/ports-and-adapters-golang/tree/master/internal/pkg/adapters/out/cartRepository
*/
type Repository interface {
	GetAllGroups(l int64, n string) ([]models.Group, string)
	CreateGroup(group models.Group) (models.Group, error)
}

//We can call this an Adapter! It connets to external service
type GroupDynamoDbRepository struct {
	client *dynamodb.DynamoDB
}

type NextKey struct {
	Id string `json:"id"`
}

var (
	tableName = aws.String(os.Getenv("GROUPS_TABLE"))
)

// Creates a DynamoDb client
func createDynamoDBClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession()) // Use aws sdk to connect to dynamoDB
	svc := dynamodb.New(sess)                  // Create DynamoDB client
	xray.AWS(svc.Client)
	return svc
}

// NewDynamoDbRepo creates a new DynamoDb Repository
func NewDynamoDbRepo() Repository {
	dbc := createDynamoDBClient()

	return &GroupDynamoDbRepository{dbc}
}

// Get all the groups users created
func (r *GroupDynamoDbRepository) GetAllGroups(limit int64, nextKey string) ([]models.Group, string) {
	// Read from DynamoDB
	var input *dynamodb.ScanInput

	if nextKey == "" {
		input = &dynamodb.ScanInput{
			TableName: tableName,
			Limit:     aws.Int64(limit),
		}
	} else {
		nk := &NextKey{}

		//We decode the key
		k, _ := url.QueryUnescape(nextKey)

		//parse the key
		json.Unmarshal([]byte(k), nk)

		input = &dynamodb.ScanInput{
			TableName: tableName,
			Limit:     aws.Int64(limit),
			//ExclusiveStartKey: nk,
			ExclusiveStartKey: map[string]*dynamodb.AttributeValue{
				"id": {
					S: aws.String(nk.Id),
				},
			},
		}
	}

	result, _ := r.client.Scan(input)

	// Construct todos from response
	var groups []models.Group
	for _, i := range result.Items {
		group := models.Group{}
		if err := dynamodbattribute.UnmarshalMap(i, &group); err != nil {
			log.Println("Failed to unmarshal")
			log.Println(err)
		}
		groups = append(groups, group)
	}

	var nxt NextKey
	if err := dynamodbattribute.UnmarshalMap(result.LastEvaluatedKey, &nxt); err != nil {
		log.Println("Failed to unmarshal")
		log.Println(err)
	}
	//log.Printf("LastEvaluatedKey: %s", nxt.Id)

	var finalKeyValue string
	if nxt.Id == "" {

		//when the next key is null it means there is no more items ot return
		finalKeyValue = string("null")
	} else {
		//fmt.Sprintf("%v", nxt.Id)
		out, _ := json.Marshal(nxt)
		finalKeyValue = string(out)
	}

	return groups, finalKeyValue
}

// Store creates a new group team in the in images table.
func (r *GroupDynamoDbRepository) CreateGroup(group models.Group) (models.Group, error) {
	// Write the new item to DynamoDB database
	item, _ := dynamodbattribute.MarshalMap(group)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: tableName,
	}

	if _, err := r.client.PutItem(input); err != nil {
		return models.Group{}, err
	}

	return group, nil
}
