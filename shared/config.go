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

// TableName provides the DynamoDB table name for this process
func TableName() string {
	r, ok := os.LookupEnv("TASKAPP_TABLE_NAME")

	if !ok {
		return DefaultTableName
	}

	return r
}
