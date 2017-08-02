package shared

import "os"

const (
	// DefaultTableName is the name of the DynamoDB table name for this app
	DefaultTableName = "taskapp"

	// DefaultAWSRegion is the default for the region for this app
	DefaultAWSRegion = "us-east-1"
)

// AWSRegion provides the aws region for this process
func AWSRegion() string {
	r, ok := os.LookupEnv("TASKAPP_REGION")

	if !ok {
		return DefaultAWSRegion
	}

	return r
}

// NotesTableName provides the DynamoDB notes table name for this process
func NotesTableName() string {
	r, ok := os.LookupEnv("TASKAPP_TASKS_TABLE_NAME")

	if !ok {
		return DefaultTableName
	}

	return r
}

// TasksTableName provides the DynamoDB tasks table name for this process
func TasksTableName() string {
	r, ok := os.LookupEnv("TASKAPP_TASKS_TABLE_NAME")

	if !ok {
		return DefaultTableName
	}

	return r
}
