package main

import (
	"encoding/json"

	apex "github.com/apex/go-apex"
)

func main() {
	// region := taskapp.AWSRegion()
	// tableName := taskapp.TasksTableName()

	// svc := dynamodb.New(session.New(&aws.Config{
	// 	Region: aws.String(region),
	// }))

	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {

		return "Success", nil
	})
}
