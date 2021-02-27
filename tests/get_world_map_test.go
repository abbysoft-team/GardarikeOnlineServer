package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWorldMap(t *testing.T) {
	TestLoginSuccessful(t)

	var request rpc.Request
	request.Data = &rpc.Request_GetWorldMapRequest{
		GetWorldMapRequest: &rpc.GetWorldMapRequest{
			Location: &rpc.IntVector2D{
				X: 2,
				Y: 2,
			},
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
	if !assert.NotNil(t, resp.GetGetWorldMapResponse(), "response isn't a get world map response") {
		return
	}

	assert.NotEmpty(t, resp.GetGetWorldMapResponse().Map.Data)
}
