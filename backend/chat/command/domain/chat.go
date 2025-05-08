package domain

import (
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

type ChatEvent struct {
	ChatId    uuid.UUID
	EventType string
	Version   int64
	Payload   []byte
	Timestamp int64
}

type ChatEventOutbox struct {
	EventId   int64
	EventType string
	Payload   []byte
	Timestamp int64
}

func (c *Chat) AsCreateEvent() ChatEvent {
	json_chat, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return ChatEvent{
		ChatId:    c.ID,
		EventType: "chat_created",
		Version:   1,
		Payload:   json_chat,
		Timestamp: c.Timestamp,
	}
}

func (e *ChatEvent) AsOutbox(eventId int64) ChatEventOutbox {
	return ChatEventOutbox{
		EventId:   eventId,
		EventType: e.EventType,
		Payload:   e.Payload,
		Timestamp: e.Timestamp,
	}
}
