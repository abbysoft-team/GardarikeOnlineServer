package model

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"fmt"
)

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

func NewError(message string, code rpc.Error) Error {
	return SimpleError{
		Message: message,
		Code:    int(code),
	}
}

func (e SimpleError) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

var ErrInternalServerError = NewError("internal server error", rpc.Error_INTERNAL_SERVER_ERROR)
var ErrInvalidUserPassword = NewError("invalid username/password combination", rpc.Error_INVALID_PASSWORD)
var ErrNotAuthorized = NewError("user not authorized", rpc.Error_NOT_AUTHORIZED)
var ErrCharacterNotFound = NewError("character not found", rpc.Error_CHARACTER_NOT_FOUND)
var ErrBadRequest = NewError("bad request", rpc.Error_BAD_REQUEST)
var ErrCharacterNotSelected = NewError("character not selected", rpc.Error_CHARACTER_NOT_SELECTED)
var ErrMessageTooLong = NewError("chat message too long", rpc.Error_MESSAGE_TOO_LONG)
var ErrUsernameIsTaken = NewError("username is already registered", rpc.Error_USERNAME_IS_ALREADY_TAKEN)
var ErrForbidden = NewError("action is forbidden", rpc.Error_FORBIDDEN)
var ErrNotEnoughResources = NewError("not enough resources", rpc.Error_NOT_ENOUGH_RESOURCES)
var ErrTownNotFound = NewError("town not found", rpc.Error_TOWN_NOT_FOUND)
