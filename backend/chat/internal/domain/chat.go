package domain

import (
	"boardroom/shared/event"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID       uuid.UUID
	SenderID uuid.UUID
	RoomID   uuid.UUID
	Message  string
	Version  int64
}

func NewChat(senderId uuid.UUID, roomId uuid.UUID, message string) *Chat {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return &Chat{
		ID:       id,
		SenderID: senderId,
		RoomID:   roomId,
		Message:  message,
		Version:  1,
	}
}

func NewEditedChat(id uuid.UUID, senderId uuid.UUID, roomId uuid.UUID, message string) *Chat {
	return &Chat{
		ID:       id,
		SenderID: senderId,
		RoomID:   roomId,
		Message:  message,
		Version:  1, // ToDo
	}
}

func (c *Chat) NewCreatedEvent() (*event.ChatEvent, error) {
	payload := event.ChatCreatedPayload{
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
	return &event.ChatEvent{
		Type:       event.EventTypeCreated,
		OccurredAt: time.Now().Unix(),
		Payload:    jsonChat,
	}, nil
}

func (c *Chat) NewUpdatedEvent() (*event.ChatEvent, error) {
	p := event.ChatEditedPayload{
		ID:      c.ID,
		RoomID:  c.RoomID,
		Message: c.Message,
	}
	jsonChat, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return &event.ChatEvent{
		Type:       event.EventTypeUpdated,
		OccurredAt: time.Now().Unix(),
		Payload:    jsonChat,
	}, nil
}

func NewDeletedEvent(roomId uuid.UUID, chatId uuid.UUID) (*event.ChatEvent, error) {
	p := event.ChatDeletedPayload{
		ID:     chatId,
		RoomID: roomId,
	}
	jsonChat, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return &event.ChatEvent{
		Type:       event.EventTypeDeleted,
		OccurredAt: time.Now().Unix(),
		Payload:    jsonChat,
	}, nil
}

//func (c *Chat) AsEditEvent() ChatEvent {
//	jsonChat, err := json.Marshal(c)
//	if err != nil {
//		panic(err)
//	}
//	return ChatEvent{
//		ChatID:    c.ID,
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
//		ChatID:    c.ID,
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
