package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func requireEvent(t *testing.T, timeout time.Duration) *rpc.Event {
	select {
	case event := <-client.eventChan:
		// Okay
		require.NotNil(t, event)
		return event
	case <-time.NewTimer(timeout).C:
		// Not okay
		t.Fatalf("No events happened after %f seconds", timeout.Seconds())
	}

	return nil
}

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

	require.NoError(t, err)

	event := requireEvent(t, time.Second)

	require.NotNil(t, event.GetChatMessageEvent())
	require.Equal(t, event.GetChatMessageEvent().Message.Type, rpc.ChatMessage_SYSTEM)
}
