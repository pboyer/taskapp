package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	apex "github.com/apex/go-apex"
	taskapp "github.com/pboyer/taskapp/shared"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type response struct {
	Items []*taskapp.Task `json:"items"`
}

type tasks []*taskapp.Task

func (s tasks) Len() int {
	return len(s)
}

func (s tasks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s tasks) Less(i, j int) bool {
	t0 := s[i]
	t1 := s[j]

	if t0.Completed == nil {
		if t1.Completed == nil {
			return *t0.Priority < *t1.Priority
		}

		return true
	}

	if t1.Completed == nil {
		return false
	}

	time0, _ := taskapp.ParseISO8601Time(*t0.Completed)
	time1, _ := taskapp.ParseISO8601Time(*t1.Completed)

	if time0.Equal(time1) {
		return *t0.Priority < *t1.Priority
	}

	return time0.After(time1)
}

func main() {
	region := taskapp.AWSRegion()
	tableName := taskapp.TasksTableName()

	svc := dynamodb.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		var task taskapp.Task

		if err := json.Unmarshal(event, &task); err != nil {
			return nil, taskapp.BadRequest(fmt.Sprintf("Failed to parse input: %v", err))
		}

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

		if len(expAttValues) == 0 {
			return nil, taskapp.BadRequest("You must supply at least one filter parameter")
		}

		if len(expAttNames) > 0 {
			input.ExpressionAttributeNames = expAttNames
		}

		result, err := svc.Scan(input)

		if err != nil {
			return nil, taskapp.InternalServerError(err)
		}

		items := make([]*taskapp.Task, len(result.Items))
		for i, v := range result.Items {
			task, err := taskapp.NewTaskFromAttributeValueMap(v)
			if err != nil {
				return nil, taskapp.InternalServerError(err)
			}
			items[i] = task
		}

		sort.Sort(tasks(items))

		return response{items}, nil
	})
}
