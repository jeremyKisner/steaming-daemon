package connector

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateSession() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String("http://dynamodb:8000"), // DynamoDB Local endpoint
		Region:   aws.String("us-west-2"),            // Change region as necessary
	})
	if err != nil {
		log.Fatal(err)
	}
	dynamoDB := dynamodb.New(sess)
	return dynamoDB
}
