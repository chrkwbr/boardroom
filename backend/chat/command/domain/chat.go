package domain

import (
	"backend/chat/event"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Chat struct {
	ID        uuid.UUID
	Sender    string
	Room      string
	Message   string
	Version   int64
	Timestamp int64
}

func (c *Chat) Edit(message string) Chat {
	return Chat{
		ID:        c.ID,
		Sender:    c.Sender,
		Room:      c.Room,
		Message:   message,
		Version:   c.Version + 1,
		Timestamp: c.Timestamp,
	}

}

type ChatEventOutbox struct {
	EventId   int64
	EventType string
	Payload   []byte
	Timestamp int64
}

func (c *Chat) AsCreateEvent() event.ChatEvent {
	jsonChat, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return event.ChatEvent{
		ChatId:    c.ID,
		EventType: event.ChatCreatedEvent,
		Version:   1,
		Payload:   jsonChat,
		Timestamp: time.Now().Unix(),
	}
}

func (c *Chat) AsEditEvent() event.ChatEvent {
	jsonChat, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return event.ChatEvent{
		ChatId:    c.ID,
		EventType: event.ChatEditedEvent,
		Version:   c.Version,
		Payload:   jsonChat,
		Timestamp: time.Now().Unix(),
	}
}

func (c *Chat) AsDeleteEvent() event.ChatEvent {
	jsonChat, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return event.ChatEvent{
		ChatId:    c.ID,
		EventType: event.ChatDeletedEvent,
		Version:   c.Version,
		Payload:   jsonChat,
		Timestamp: time.Now().Unix(),
	}
}

func AsOutbox(eventId int64, e event.ChatEvent) ChatEventOutbox {
	marshaledEvent, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return ChatEventOutbox{
		EventId:   eventId,
		EventType: e.EventType,
		Payload:   marshaledEvent,
		Timestamp: e.Timestamp,
	}
}
