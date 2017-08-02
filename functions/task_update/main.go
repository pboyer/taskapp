package main

import (
	apex "github.com/apex/go-apex"
	taskapp "github.com/pboyer/taskapp/shared"
)

func main() {
	apex.HandleFunc(taskapp.TaskPutFunc(false))
}
