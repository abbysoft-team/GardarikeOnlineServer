package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPlaceTown(t *testing.T) {
	TestSelectCharacter(t)

	var request rpc.Request
	request.Data = &rpc.Request_PlaceTownRequest{
		PlaceTownRequest: &rpc.PlaceTownRequest{
			SessionID: sessionID,
			Location:  nil,
			Name:      fmt.Sprintf("test-town-%d", time.Now().Unix()),
		},
	}

	resp, err := client.SendRequest(request)

	if !assert.NoError(t, err, "request error is not nil") {
		return
	}
	if !assert.NotNil(t, resp, "response is nil") {
		return
	}
	if !assert.NotNil(t, resp.GetPlaceTownResponse(), "response isn't a place town response") {
		return
	}

	assert.NotNil(t, resp.GetPlaceTownResponse().Location)
}
