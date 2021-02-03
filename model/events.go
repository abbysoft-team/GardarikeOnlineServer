package model

import (
	"abbysoft/gardarike-online/model/consts"
	rpc "abbysoft/gardarike-online/rpc/generated"
)

func NewChatMessageEvent(message ChatMessage) EventWrapper {
	return EventWrapper{
		Event: &rpc.Event{
			Payload: &rpc.Event_ChatMessageEvent{
				ChatMessageEvent: &rpc.NewChatMessageEvent{
					Message: message.ToRPC(),
				},
			},
		},
		Topic: consts.GlobalTopic,
	}
}

func NewSystemChatMessageEvent(text string) EventWrapper {
	return NewChatMessageEvent(ChatMessage{
		ID:       0,
		Sender:   consts.SystemUserName,
		Text:     text,
		IsSystem: true,
	})
}
