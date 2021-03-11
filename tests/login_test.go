// +build !remote_tests

package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err, "Error while making request")

	loginResponse := response.GetLoginResponse()
	require.NotNil(t, loginResponse, "Response is nil")
	require.NotEmpty(t, loginResponse.SessionID)

	sessionID = loginResponse.SessionID
}
