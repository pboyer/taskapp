package main

import (
	"context"
	"encoding/json"
	"time"

	apex "github.com/apex/go-apex"
	taskapp "github.com/pboyer/taskapp/shared"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		sess := session.New(&aws.Config{Region: aws.String(taskapp.DefaultAWSRegion)})
		svc := dynamodb.New(sess)

		// timeout policy?
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tableName := "taskapp"
		keyCond := "#user = :u"
		user := "foobar"
		userFieldName := "user"

		return svc.QueryWithContext(ctx, &dynamodb.QueryInput{
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
