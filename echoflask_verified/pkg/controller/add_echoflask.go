package controller

import (
	"echoflask/pkg/controller/echoflask"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, echoflask.Add)
}
