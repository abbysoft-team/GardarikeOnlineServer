// +build !remote_tests

package tests

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateAccount(t *testing.T) {
	var request rpc.Request
	request.Data = &rpc.Request_CreateAccountRequest{
		CreateAccountRequest: &rpc.CreateAccountRequest{
			Login:    fmt.Sprintf("test%d", time.Now().Unix()),
			Password: "test",
		},
	}

	resp, err := client.SendRequest(request)
	if !assert.NoError(t, err, "request error is not nil") {
		return
	}

	if !assert.NotNil(t, resp, "response is nil") {
		return
	}

	if !assert.NotNil(t, resp.GetCreateAccountResponse(), "response isn't a create account response") {
		return
	}

	assert.NotZero(t, resp.GetCreateAccountResponse().Id, "account ID is zero")
}

func TestCreateAccount_AlreadyExists(t *testing.T) {
	var request rpc.Request
	request.Data = &rpc.Request_CreateAccountRequest{
		CreateAccountRequest: &rpc.CreateAccountRequest{
			Login:    "test",
			Password: "test",
		},
	}

	_, err := client.SendRequest(request)

	if !assert.EqualError(t, err, model.ErrUsernameIsTaken.Error(),
		"error isn't a model.ErrUsernameIsTaken error") {
		return
	}
}

func TestCreateAccount_EmptyUsername(t *testing.T) {
	var request rpc.Request
	request.Data = &rpc.Request_CreateAccountRequest{
		CreateAccountRequest: &rpc.CreateAccountRequest{
			Login:    "test",
			Password: "test",
		},
	}

	_, err := client.SendRequest(request)

	assert.NotNil(t, err, "error isn't nil")
}
