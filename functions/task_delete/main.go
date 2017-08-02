package main

import (
	"encoding/json"
	"fmt"

	taskapp "github.com/pboyer/taskapp/shared"

	apex "github.com/apex/go-apex"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type request struct {
	ID *string `json:"id"`
}

func main() {
	region := taskapp.AWSRegion()
	tableName := taskapp.TasksTableName()

	svc := dynamodb.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		var req request

		if err := json.Unmarshal(event, &req); err != nil {
			return nil, taskapp.BadRequest(fmt.Sprintf("%v", err))
		}

		if req.ID == nil {
			return nil, taskapp.BadRequest("Deletion requires a task 'id'")
		}

		_, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: &tableName,
			Key: map[string]*dynamodb.AttributeValue{
				"id": &dynamodb.AttributeValue{S: req.ID},
			},
		})

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeResourceNotFoundException:
					return nil, taskapp.BadRequest("That id does not exist")
				default:
					return nil, taskapp.InternalServerError(err)
				}
			}
		}

		return "Success", nil
	})
}
