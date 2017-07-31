package main

import (
	"encoding/json"
	"errors"

	taskapp "github.com/pboyer/taskapp/shared"

	apex "github.com/apex/go-apex"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type request struct {
	ID *string `json:"id"`
}

func main() {
	svc := dynamodb.New(session.New(&aws.Config{
		Region: aws.String(taskapp.DefaultAWSRegion),
	}))

	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		var req request

		if err := json.Unmarshal(event, &req); err != nil {
			return nil, err
		}

		if req.ID == nil {
			return nil, errors.New("Deletion requires a task 'id'")
		}

		tableName := taskapp.DefaultTableName
		return svc.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: &tableName,
			Key: map[string]*dynamodb.AttributeValue{
				"id": &dynamodb.AttributeValue{S: req.ID},
			},
		})
	})
}
