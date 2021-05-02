package groupsAccess

import (
	"log"
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
	GetAllGroups() []models.Group
	CreateGroup(group models.Group) (models.Group, error)
}

//We can call this an Adapter! It connets to external service
type GroupDynamoDbRepository struct {
	client *dynamodb.DynamoDB
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
func (r *GroupDynamoDbRepository) GetAllGroups() []models.Group {
	// Read from DynamoDB
	input := &dynamodb.ScanInput{
		TableName: tableName,
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

	return groups
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
