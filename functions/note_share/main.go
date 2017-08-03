package main

import (
	"encoding/json"
	"fmt"

	taskapp "github.com/pboyer/taskapp/shared"

	apex "github.com/apex/go-apex"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type request struct {
	ID              *string `json:"id"`
	OriginatingUser *string `json:"originating_user"` // this user sent the request, must be the creator of the note
	NewCollaborator *string `json:"new_collaborator"` // the is the user being added, must not duplicate existing collaboarators
}

func main() {
	region := taskapp.AWSRegion()
	tableName := taskapp.NotesTableName()

	svc := dynamodb.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		var req request

		if err := json.Unmarshal(event, &req); err != nil {
			return nil, taskapp.BadRequest(fmt.Sprintf("%v", err))
		}

		if req.ID == nil {
			return nil, taskapp.BadRequest("Requires a note 'id'")
		}

		if req.OriginatingUser == nil {
			return nil, taskapp.BadRequest("Requires an 'originating_user' field")
		}

		if err := taskapp.ValidateEmailString(*req.OriginatingUser); err != nil {
			return nil, taskapp.BadRequest("Invalid email in 'originating_user' field")
		}

		if req.NewCollaborator == nil {
			return nil, taskapp.BadRequest("Requires an 'new_collaborator' field")
		}

		if err := taskapp.ValidateEmailString(*req.NewCollaborator); err != nil {
			return nil, taskapp.BadRequest("Invalid email in 'new_collaborator' field")
		}

		item, err := svc.GetItem(&dynamodb.GetItemInput{
			TableName: &tableName,
			Key: map[string]*dynamodb.AttributeValue{
				"id": &dynamodb.AttributeValue{S: req.ID},
			},
		})

		if err != nil {
			return nil, taskapp.InternalServerError(err)
		}

		if item.Item == nil {
			return nil, taskapp.BadRequest("That id does not exist")
		}

		note, err := taskapp.NewNoteFromAttributeValueMap(item.Item)
		if err != nil {
			return nil, taskapp.InternalServerError(err)
		}

		if *note.Creator != *req.OriginatingUser {
			return nil, taskapp.BadRequest("The 'originating_user' field does not match the note 'creator'.")
		}

		if *note.Creator == *req.NewCollaborator {
			return nil, taskapp.BadRequest("The 'creator' cannot be made a collaborator.")
		}

		if note.Collaborators == nil {
			note.Collaborators = make([]*string, 0, 1)
		} else {
			for _, c := range note.Collaborators {
				if *c == *req.NewCollaborator {
					// duplicate
					return nil, taskapp.BadRequest("The 'originating_user' is already a collaborator.")
				}
			}
		}

		note.Collaborators = append(note.Collaborators, req.NewCollaborator)

		_, err = svc.PutItem(&dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      note.ToAttributeValueMap(),
		})

		if err != nil {
			return nil, taskapp.InternalServerError(err)
		}

		return note, nil
	})
}
