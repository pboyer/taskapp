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

// NotePutFunc provides a function that adds or updates a note based on the event JSON. Setting generateNewID to true
// causes a new Note ID to be generated, appropriate for insertion of a new note into the db
func NotePutFunc(generateNewID bool) func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
	region := AWSRegion()
	tableName := NotesTableName()

	svc := dynamodb.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	return func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		var note Note

		if err := json.Unmarshal(event, &note); err != nil {
			return nil, BadRequest(fmt.Sprintf("Failed to parse input: %v", err))
		}

		if generateNewID {
			newid := uuid.NewV4().String()
			note.ID = &newid
		}

		if err := note.Validate(); err != nil {
			return nil, BadRequest(err)
		}

		_, err := svc.PutItem(&dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      note.ToAttributeValueMap(),
		})

		if err != nil {
			// TODO log for reconaissance, give unique error code
			return nil, InternalServerError(err)
		}

		return note, nil
	}
}
