package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetResources(t *testing.T) {
	TestSelectCharacter(t)

	var request rpc.Request
	request.Data = &rpc.Request_GetResourcesRequest{
		GetResourcesRequest: &rpc.GetResourcesRequest{
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
	if !assert.NotNil(t, resp.GetGetResourcesResponse(), "response isn't a get resources response") {
		return
	}
}
