package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectCharacter(t *testing.T) {
	TestLoginSuccessful(t)

	var request rpc.Request
	request.Data = &rpc.Request_SelectCharacterRequest{
		SelectCharacterRequest: &rpc.SelectCharacterRequest{
			CharacterID: 5,
			SessionID:   sessionID,
		},
	}

	resp, err := client.SendRequest(request)

	if !assert.NoError(t, err, "request error is not nil") {
		return
	}
	if !assert.NotNil(t, resp, "response is nil") {
		return
	}
	if !assert.NotNil(t, resp.GetSelectCharacterResponse(), "response isn't a select character response") {
		return
	}
}
