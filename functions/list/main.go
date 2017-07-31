package main

import (
	"encoding/json"
	"fmt"
	"strings"

	apex "github.com/apex/go-apex"
	taskapp "github.com/pboyer/taskapp/shared"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	svc := dynamodb.New(session.New(&aws.Config{
		Region: aws.String(taskapp.DefaultAWSRegion),
	}))

	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		var task taskapp.Task

		if err := json.Unmarshal(event, &task); err != nil {
			return nil, fmt.Errorf("Failed to parse input: %v", err)
		}

		tableName := taskapp.DefaultTableName

		keyConds := []string{}
		expAttNames := map[string]*string{}
		expAttValues := map[string]*dynamodb.AttributeValue{}

		if task.ID != nil {
			keyConds = append(keyConds, "id = :i")
			expAttValues[":i"] = &dynamodb.AttributeValue{S: task.ID}
		}

		if task.User != nil {
			keyConds = append(keyConds, "#user = :u")
			expAttValues[":u"] = &dynamodb.AttributeValue{S: task.User}

			// must be remapped as it is a reserved field name
			userFieldName := "user"
			expAttNames["#user"] = &userFieldName
		}

		if task.Description != nil {
			keyConds = append(keyConds, "description = :d")
			expAttValues[":d"] = &dynamodb.AttributeValue{S: task.Description}
		}

		if task.Priority != nil {
			priorityStr := fmt.Sprintf("%d", *task.Priority)
			keyConds = append(keyConds, "priority = :p")
			expAttValues[":p"] = &dynamodb.AttributeValue{N: &priorityStr}
		}

		if task.Completed != nil {
			keyConds = append(keyConds, "completed = :c")
			expAttValues[":c"] = &dynamodb.AttributeValue{S: task.Completed}
		}

		keyCond := strings.Join(keyConds, " and ")

		input := &dynamodb.ScanInput{
			TableName:                 &tableName,
			FilterExpression:          &keyCond,
			ExpressionAttributeValues: expAttValues,
		}

		if len(expAttNames) > 0 {
			input.ExpressionAttributeNames = expAttNames
		}

		result, err := svc.Scan(input)

		if err != nil {
			// TODO log for reconaissance, give unique error code
			return nil, fmt.Errorf("Internal error: %v", err)
		}

		items := make([]*taskapp.Task, len(result.Items))
		for i, v := range result.Items {
			task, err := taskapp.NewTaskFromAttributeValueMap(v)
			if err != nil {
				// TODO log for reconaissance, give unique error code
				return nil, fmt.Errorf("Internal error: %v", err)
			}
			items[i] = task
		}

		return items, nil
	})
}
