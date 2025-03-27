package models

import "errors"

var (
	ErrNotFound     = errors.New("the item not found")
	ErrExists       = errors.New("the item already exists")
	ErrNotEnough    = errors.New("product ID is missing")
	SuccesMsg       = string("Item successfully added")
	SuccesDeleteMsg = string("Item successfully deleted")
)

type ErrStruct struct {
	ErrMsg string `json:"error"`
}
