package domain

import (
	"backend/chat/event"
	"encoding/json"
	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID
	Sender    string
	Room      string
	Message   string
	Timestamp int64
}

type ChatEventOutbox struct {
	EventId   int64
	EventType string
	Payload   []byte
	Timestamp int64
}

func (c *Chat) AsCreateEvent() event.ChatEvent {
	json_chat, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return event.ChatEvent{
		ChatId:    c.ID,
		EventType: "chat_created",
		Version:   1,
		Payload:   json_chat,
		Timestamp: c.Timestamp,
	}
}

func AsOutbox(eventId int64, e event.ChatEvent) ChatEventOutbox {
	return ChatEventOutbox{
		EventId:   eventId,
		EventType: e.EventType,
		Payload:   e.Payload,
		Timestamp: e.Timestamp,
	}
}
