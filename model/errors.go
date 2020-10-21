package model

import "fmt"

type Error interface {
	error
	GetMessage() string
	GetCode() int
}

type SimpleError struct {
	Message string
	Code    int
}

func (e SimpleError) GetMessage() string {
	return e.Message
}

func (e SimpleError) GetCode() int {
	return e.Code
}

func NewError(message string, code int) Error {
	return SimpleError{
		Message: message,
		Code:    code,
	}
}

func (e SimpleError) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

var ErrInternalServerError = NewError("internal server error", 1)
var ErrInvalidUserPassword = NewError("invalid username/password combination", 2)
var ErrNotAuthorized = NewError("user not authorized", 3)
var ErrCharacterNotFound = NewError("character not found", 4)
var ErrBuildingNotFound = NewError("building not found", 5)
var ErrNoEnoughMoney = NewError("no enough money", 6)
var ErrBuildingSpotIsBusy = NewError("building spot is busy", 6)
