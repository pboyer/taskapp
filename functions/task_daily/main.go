package main

import (
	"encoding/json"
	"fmt"
	"os"

	apex "github.com/apex/go-apex"
)

func main() {
	// region := taskapp.AWSRegion()
	// tableName := taskapp.TasksTableName()

	// svc := dynamodb.New(session.New(&aws.Config{
	// 	Region: aws.String(region),
	// }))

	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		fmt.Fprintf(os.Stderr, "Send emails!")

		return "Success", nil
	})
}
