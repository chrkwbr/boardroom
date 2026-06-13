package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID
	SenderID  uuid.UUID
	RoomID    uuid.UUID
	Message   string
	Version   int64
	Timestamp int64
}

func NewChat(senderId uuid.UUID, roomId uuid.UUID, message string) *Chat {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return &Chat{
		ID:        id,
		SenderID:  senderId,
		RoomID:    roomId,
		Message:   message,
		Version:   1,
		Timestamp: time.Now().Unix(),
	}
}

func NewEditedChat(id uuid.UUID, senderId uuid.UUID, roomId uuid.UUID, message string) *Chat {
	return &Chat{
		ID:        id,
		SenderID:  senderId,
		RoomID:    roomId,
		Message:   message,
		Version:   1,
		Timestamp: time.Now().Unix(),
	}
}

//func (c *Chat) Edit(message string) Chat {
//	return Chat{
//		ID:        c.ID,
//		Sender:    c.Sender,
//		Room:      c.Room,
//		Message:   message,
//		Version:   c.Version + 1,
//		Timestamp: c.Timestamp,
//	}
//
//}

//type ChatEventOutbox struct {
//	EventId   int64
//	EventType string
//	Payload   []byte
//	Timestamp int64
//}

func (c *Chat) NewCreatedEvent() *ChatEvent {
	payload := ChatCreatedPayload{
		ID:       c.ID,
		RoomID:   c.RoomID,
		SenderID: c.SenderID,
		Message:  c.Message,
		Version:  c.Version,
	}
	jsonChat, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return &ChatEvent{
		Type:       EventTypeCreated,
		OccurredAt: c.Timestamp,
		Payload:    jsonChat,
	}
}

func (c *Chat) NewUpdatedEvent() *ChatEvent {
	p := ChatEditedPayload{
		ID:      c.ID,
		RoomID:  c.RoomID,
		Message: c.Message,
	}
	jsonChat, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return &ChatEvent{
		Type:       EventTypeUpdated,
		OccurredAt: c.Timestamp,
		Payload:    jsonChat,
	}
}

//func (c *Chat) AsEditEvent() ChatEvent {
//	jsonChat, err := json.Marshal(c)
//	if err != nil {
//		panic(err)
//	}
//	return ChatEvent{
//		ChatId:    c.ID,
//		EventType: ChatEditedPayload,
//		Version:   c.Version,
//		Payload:   jsonChat,
//		Timestamp: time.Now().Unix(),
//	}
//}
//
//func (c *Chat) AsDeleteEvent() ChatEvent {
//	jsonChat, err := json.Marshal(c)
//	if err != nil {
//		panic(err)
//	}
//	return ChatEvent{
//		ChatId:    c.ID,
//		EventType: ChatDeletedEvent,
//		Version:   c.Version,
//		Payload:   jsonChat,
//		Timestamp: time.Now().Unix(),
//	}
//}
//
//func AsOutbox(eventId int64, e ChatEvent) ChatEventOutbox {
//	marshaledEvent, err := json.Marshal(e)
//	if err != nil {
//		panic(err)
//	}
//	return ChatEventOutbox{
//		EventId:   eventId,
//		EventType: e.EventType,
//		Payload:   marshaledEvent,
//		Timestamp: e.Timestamp,
//	}
//}
