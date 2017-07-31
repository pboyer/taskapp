package main

import (
	"encoding/json"

	apex "github.com/apex/go-apex"
	taskapp "github.com/pboyer/taskapp/shared"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		svc := dynamodb.New(session.New(&aws.Config{
			Region: aws.String(taskapp.DefaultAWSRegion),
		}))

		tableName := taskapp.DefaultTableName
		keyCond := "#user = :u"
		user := "peter."
		userFieldName := "user"

		return svc.Query(&dynamodb.QueryInput{
			TableName:              &tableName,
			KeyConditionExpression: &keyCond,
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":u": &dynamodb.AttributeValue{
					S: &user,
				},
			},
			ExpressionAttributeNames: map[string]*string{
				"#user": &userFieldName,
			},
		})
	})
}
