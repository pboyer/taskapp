package main

import (
	"encoding/json"

	taskapp "github.com/pboyer/taskapp/shared"

	apex "github.com/apex/go-apex"
	uuid "github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		var t taskapp.Task

		if err := json.Unmarshal(event, &t); err != nil {
			return nil, err
		}

		newid := uuid.NewV4().String()
		t.ID = &newid

		if err := t.Validate(); err != nil {
			return nil, err
		}

		svc := dynamodb.New(session.New(&aws.Config{
			Region: aws.String(taskapp.DefaultAWSRegion),
		}))

		tableName := taskapp.DefaultTableName
		return svc.PutItem(&dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      t.ToAttributeValueMap(),
		})
	})
}
