package shared

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Note represents a simple note sharable between users
type Note struct {
	ID            *string   `json:"id"`      // DynamoDB partition key
	Creator       *string   `json:"creator"` // required
	Text          *string   `json:"text"`    // required
	Collaborators []*string `json:"collaborators"`
}

// Validate validates the Note according to its schema.
func (n *Note) Validate() error {
	if err := n.validateID(); err != nil {
		return fmt.Errorf("The 'id' attribute is invalid: %v", err)
	}

	if err := n.validateCreator(); err != nil {
		return fmt.Errorf("The 'creator' attribute is invalid: %v", err)
	}

	if err := n.validateText(); err != nil {
		return fmt.Errorf("The 'text' attribute is invalid: %v", err)
	}

	if err := n.validateCollaborators(); err != nil {
		return fmt.Errorf("The 'collabarators' attribute is invalid: %v", err)
	}

	return nil
}

func NewNoteFromAttributeValueMap(m map[string]*dynamodb.AttributeValue) (*Note, error) {
	note := &Note{}

	if id, ok := m["id"]; ok {
		note.ID = id.S
	}

	if creator, ok := m["creator"]; ok {
		note.Creator = creator.S
	}

	if text, ok := m["text"]; ok {
		note.Text = text.S
	}

	if collaborators, ok := m["collaborators"]; ok {
		note.Collaborators = collaborators.SS
	}

	// This validate step should never fail as its decoding data already entered into the DB.
	// We do it anyways to avoid serving up invalid content in the event of failure. It could be
	// turned off in future builds if performance justifies it.
	if err := note.Validate(); err != nil {
		return nil, err
	}

	return note, nil
}

// ToAttributeValueMap turns the Task as an AttributeValue map for use in Amazon DynamoDB API's. This function
// does not validate the Note nor does it assume validation was completed.
func (n *Note) ToAttributeValueMap() map[string]*dynamodb.AttributeValue {
	res := map[string]*dynamodb.AttributeValue{}

	if n.ID != nil {
		res["id"] = &dynamodb.AttributeValue{S: n.ID}
	}

	if n.Creator != nil {
		res["creator"] = &dynamodb.AttributeValue{S: n.Creator}
	}

	if n.Text != nil {
		res["text"] = &dynamodb.AttributeValue{S: n.Text}
	}

	if n.Collaborators != nil {
		res["collaborators"] = &dynamodb.AttributeValue{SS: n.Collaborators}
	}

	return res
}

func (n *Note) validateID() error {
	if n.ID == nil {
		return errors.New("Attribute is required")
	}

	return validateIDString(*n.ID)
}

func (n *Note) validateCreator() error {
	if n.Creator == nil {
		return errors.New("Attribute is required")
	}

	return ValidateEmailString(*n.Creator)
}

func (n *Note) validateText() error {
	if n.Text == nil {
		return errors.New("Attribute is required")
	}

	desc := *n.Text

	if len(desc) < 1 {
		return errors.New("The attribute must be at least 1 character long")
	}

	return nil
}

func (n *Note) validateCollaborators() error {
	return nil
}
