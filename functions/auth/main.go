package main

import (
	"encoding/json"
	"fmt"

	apex "github.com/apex/go-apex"
)

type request struct {
	authorizationToken *string `json:"authorizationToken"`
	methodArn          *string `json:"methodArn"`
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, actx *apex.Context) (interface{}, error) {
		var req request

		if err := json.Unmarshal(event, &req); err != nil {
			return nil, err
		}

		// lookup username
		return nil, fmt.Errorf("Not implemented. Input: %v", req)
	})
}
