package shared

import (
	"encoding/json"
	"fmt"

	apex "github.com/apex/go-apex"
	uuid "github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// TaskPutFunc provides a function that adds or updates a task based on the event JSON. Setting generateNewID to true
// causes a new Task ID to be generated, appropriate for insertion of a new task into the db
func TaskPutFunc(generateNewID bool) func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
	region := AWSRegion()
	tableName := TasksTableName()

	svc := dynamodb.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	return func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		var task Task

		if err := json.Unmarshal(event, &task); err != nil {
			return nil, BadRequest(fmt.Sprintf("Failed to parse input: %v", err))
		}

		if generateNewID {
			newid := uuid.NewV4().String()
			task.ID = &newid
		}

		if err := task.Validate(); err != nil {
			return nil, BadRequest(err)
		}

		_, err := svc.PutItem(&dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      task.ToAttributeValueMap(),
		})

		if err != nil {
			// TODO log for reconaissance, give unique error code
			return nil, InternalServerError(err)
		}

		return task, nil
	}
}
