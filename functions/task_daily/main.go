package main

import (
	"encoding/json"
	"fmt"
	"os"

	apex "github.com/apex/go-apex"
	taskapp "github.com/pboyer/taskapp/shared"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	region := taskapp.AWSRegion()
	tableName := taskapp.TasksTableName()

	svc := dynamodb.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		fmt.Fprintf(os.Stderr, "Sending emails")

		keyCond := "completed = :c"
		s := ""
		expAttValues := map[string]*dynamodb.AttributeValue{
			":c": &dynamodb.AttributeValue{S: &s},
		}

		input := &dynamodb.ScanInput{
			TableName:                 &tableName,
			FilterExpression:          &keyCond,
			ExpressionAttributeValues: expAttValues,
		}

		result, err := svc.Scan(input)

		if err != nil {
			return nil, taskapp.InternalServerError(err)
		}

		// collect the task lists to email
		userTaskLists := map[string][]*taskapp.Task{}
		for _, v := range result.Items {
			task, err := taskapp.NewTaskFromAttributeValueMap(v)
			if err != nil {
				return nil, taskapp.InternalServerError(err)
			}

			if task.User == nil {
				continue
			}

			list, ok := userTaskLists[*task.User]
			if !ok {
				list = make([]*taskapp.Task, 0, 1)
			}

			userTaskLists[*task.User] = append(list, task)
		}

		// We do NOT send the emails here. This could be added once proper authentication, authorization,
		// subscription, and unsubscription is implemented.
		for u, list := range userTaskLists {
			fmt.Fprintf(os.Stderr, "Sending emails to %s: %v", u, list)
		}

		return "Success", nil
	})
}
