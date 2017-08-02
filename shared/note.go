package shared

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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
		return nil
	}

	return validateEmailString(*n.Creator)
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
