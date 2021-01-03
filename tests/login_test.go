// +build !remote_tests

package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoginSuccessful(t *testing.T) {
	var request rpc.Request
	request.Data = &rpc.Request_LoginRequest{
		LoginRequest: &rpc.LoginRequest{
			Username: "test",
			Password: "test",
		},
	}

	response, err := client.SendRequest(request)
	if !assert.NoError(t, err, "Error while making request") {
		return
	}

	loginResponse := response.GetLoginResponse()
	if !assert.NotNil(t, loginResponse, "Response is nil") {
		return
	}

	if !assert.NotEmpty(t, loginResponse.SessionID) {
		return
	}

	sessionID = loginResponse.SessionID
}
