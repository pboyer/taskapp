package main

import (
	"encoding/json"

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

		svc := dynamodb.New(session.New(&aws.Config{
			Region: aws.String(taskapp.DefaultAWSRegion),
		}))

		tableName := taskapp.TableName
		completedDefault := "0"

		return svc.PutItem(&dynamodb.PutItemInput{
			TableName: &tableName,
			Item: map[string]*dynamodb.AttributeValue{
				"user":      &dynamodb.AttributeValue{S: m.User},
				"completed": &dynamodb.AttributeValue{S: &completedDefault},
			},
		})
	})
}
