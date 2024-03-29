package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Task is a simple thing to do. See swagger.json for a full schema.
type Task struct {
	ID          *string `json:"id"` // DynamoDB partition key
	User        *string `json:"user"`
	Description *string `json:"description"` // required
	Priority    *uint32 `json:"priority"`    // required
	Completed   *string `json:"completed"`
}

// Validate validates the Task according to its schema.
func (t *Task) Validate() error {
	if err := t.validateID(); err != nil {
		return fmt.Errorf("The 'id' attribute is invalid: %v", err)
	}

	if err := t.validateUser(); err != nil {
		return fmt.Errorf("The 'user' attribute is invalid: %v", err)
	}

	if err := t.validateDescription(); err != nil {
		return fmt.Errorf("The 'description' attribute is invalid: %v", err)
	}

	if err := t.validatePriority(); err != nil {
		return fmt.Errorf("The 'priority' attribute is invalid: %v", err)
	}

	if err := t.validateCompleted(); err != nil {
		return fmt.Errorf("The 'completed' attribute is invalid: %v", err)
	}

	return nil
}

func NewTaskFromAttributeValueMap(m map[string]*dynamodb.AttributeValue) (*Task, error) {
	task := &Task{}

	if id, ok := m["id"]; ok {
		task.ID = id.S
	}

	if user, ok := m["user"]; ok {
		task.User = user.S
	}

	if description, ok := m["description"]; ok {
		task.Description = description.S
	}

	if priority, ok := m["priority"]; ok {
		num, err := strconv.ParseUint(*priority.N, 10, 32)

		if err != nil {
			return nil, fmt.Errorf("Unexpected error parsing priority from database: %v", err)
		}

		num32 := uint32(num)
		task.Priority = &num32
	}

	if completed, ok := m["completed"]; ok {
		task.Completed = completed.S
	}

	// This validate step should never fail as its decoding data already entered into the DB.
	// We do it anyways to avoid serving up invalid content in the event of failure. It could be
	// turned off in future builds if performance justifies it.
	if err := task.Validate(); err != nil {
		return nil, err
	}

	return task, nil
}

// ToAttributeValueMap turns the Task as an AttributeValue map for use in Amazon DynamoDB API's. This function
// does not validate the Task nor does it assume validation was completed.
func (t *Task) ToAttributeValueMap() map[string]*dynamodb.AttributeValue {
	res := map[string]*dynamodb.AttributeValue{}

	if t.ID != nil {
		res["id"] = &dynamodb.AttributeValue{S: t.ID}
	}

	if t.User != nil {
		res["user"] = &dynamodb.AttributeValue{S: t.User}
	}

	if t.Description != nil {
		res["description"] = &dynamodb.AttributeValue{S: t.Description}
	}

	if t.Priority != nil {
		// must be transmitted as a number
		priorityStr := fmt.Sprintf("%d", *t.Priority)
		res["priority"] = &dynamodb.AttributeValue{N: &priorityStr}
	}

	if t.Completed != nil {
		res["completed"] = &dynamodb.AttributeValue{S: t.Completed}
	}

	return res
}

var (
	// emailRegexp is unashamedly copied from https://github.com/badoux/checkmail/blob/master/checkmail.go
	// Go proverb: A little copying is better than a little dependency.
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func (t *Task) validateID() error {
	if t.ID == nil {
		return errors.New("Attribute is required")
	}

	return validateIDString(*t.ID)
}

func validateIDString(s string) error {
	if _, err := uuid.FromString(s); err != nil {
		return err
	}

	return nil
}

func (t *Task) validateUser() error {
	if t.User == nil {
		return nil
	}

	return ValidateEmailString(*t.User)
}

func ValidateEmailString(s string) error {
	if l := len(s); l < 5 || l > 254 {
		return errors.New("Must be between 9 and 254 characters")
	}

	if !emailRegexp.MatchString(s) {
		return errors.New("Improperly formatted email")
	}

	return nil
}

func (t *Task) validatePriority() error {
	if t.Priority == nil {
		return errors.New("Attribute is required")
	}

	if *t.Priority > 9 {
		return errors.New("Attribute must be between 0 and 9")
	}

	return nil
}

func (t *Task) validateDescription() error {
	if t.Description == nil {
		return errors.New("Attribute is required")
	}

	desc := *t.Description

	if len(desc) < 1 {
		return errors.New("The attribute must be at least 1 character long")
	}

	return nil
}

// ParseISO8601Time parses a string in that format into a time.Time struct
func ParseISO8601Time(t string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05-07:00", t)
}

func (t *Task) validateCompleted() error {
	if t.Completed == nil {
		return nil
	}

	if *t.Completed == "" {
		return nil
	}

	_, err := ParseISO8601Time(*t.Completed)
	if err == nil {
		return err
	}

	return fmt.Errorf("Not a valid ISO8601 date: %v", err)
}

// MarshalJSON marshals a task to json, setting the completed field to "" when empty
// as expected by the examples
func (t *Task) MarshalJSON() ([]byte, error) {
	completed := t.Completed
	if completed == nil {
		c := ""
		completed = &c
	}

	return json.Marshal(&struct {
		ID          *string `json:"id"`
		User        *string `json:"user"`
		Description *string `json:"description"`
		Priority    *uint32 `json:"priority"`
		Completed   *string `json:"completed"`
	}{
		ID:          t.ID,
		User:        t.User,
		Description: t.Description,
		Priority:    t.Priority,
		Completed:   completed,
	})
}
