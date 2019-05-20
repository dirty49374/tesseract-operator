package controller

import (
	"github.com/dirty49374/tesseract-operator/pkg/controller/outgoingportal"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, outgoingportal.Add)
}
