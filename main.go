package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"strings"
)

type Item struct {
	EventId     string `json:"eventId"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var list []string

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Get back the title, year, and rating
	proj := expression.NamesList(expression.Name("eventId"), expression.Name("description"), expression.Name("title"))

	expr, err := expression.NewBuilder().WithProjection(proj).Build()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("RichmondEvents"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	for _, i := range result.Items {
		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		itemJson, err := json.Marshal(item)

		fmt.Println("Item: ", string(itemJson))

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		list = append(list, string(itemJson))
	}

	resultString := "[" + strings.Join(list, ",") + "]"

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       resultString,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}

func main() {
	lambda.Start(Handler)
}
