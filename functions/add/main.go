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

type request struct {
	User        *string `json:"user"`
	Description *string `json:"description"` // required
	Priority    *uint32 `json:"priority"`    // required
	Completed   *string `json:"completed"`
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		var m request

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}

		sess := session.New(&aws.Config{Region: aws.String(taskapp.DefaultAWSRegion)})
		svc := dynamodb.New(sess)

		// timeout policy?
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tableName := taskapp.TableName
		completedDefault := "0"

		return svc.PutItemWithContext(ctx, &dynamodb.PutItemInput{
			TableName: &tableName,
			Item: map[string]*dynamodb.AttributeValue{
				"user":      &dynamodb.AttributeValue{S: m.User},
				"completed": &dynamodb.AttributeValue{S: &completedDefault},
			},
		})
	})
}
