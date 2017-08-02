package main

import (
	taskapp "github.com/pboyer/taskapp/shared"

	apex "github.com/apex/go-apex"
)

func main() {
	apex.HandleFunc(taskapp.TaskPutFunc(true))
}
