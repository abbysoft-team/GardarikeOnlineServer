// +build !remote_tests

package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateCharacter(t *testing.T) {
	TestLoginSuccessful(t)

	var request rpc.Request
	request.Data = &rpc.Request_CreateCharacterRequest{
		CreateCharacterRequest: &rpc.CreateCharacterRequest{
			Name:      fmt.Sprintf("test%d", time.Now().Unix()),
			SessionID: sessionID,
		},
	}

	resp, err := client.SendRequest(request)
	if !assert.NoError(t, err, "request error is not nil") {
		return
	}

	if !assert.NotNil(t, resp, "response is nil") {
		return
	}

	if !assert.NotNil(t, resp.GetCreateCharacterResponse(), "response isn't a create character response") {
		return
	}
	assert.NotZero(t, resp.GetCreateCharacterResponse().Id, "character ID is zero")
}
