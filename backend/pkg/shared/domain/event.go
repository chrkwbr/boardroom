package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeCreated EventType = "created"
	EventTypeUpdated EventType = "updated"
	EventTypeDeleted EventType = "deleted"
)

type ChatEvent struct {
	Type       EventType `json:"type"`
	OccurredAt int64     `json:"occurred_at"`
	Payload    []byte    `json:"payload"`
}

type ChatCreatedPayload struct {
	ID       uuid.UUID `json:"id"`
	RoomID   uuid.UUID `json:"room_id"`
	SenderID uuid.UUID `json:"sender_id"`
	Message  string    `json:"message"`
	Version  int64     `json:"version"`
}

type ChatEditedPayload struct {
	ID      uuid.UUID `json:"id"`
	RoomID  uuid.UUID `json:"room_id"`
	Message string    `json:"message"`
}

type ChatDeletedPayload struct {
	ID     uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"room_id"`
}

func NewCreatedEvent(c *Chat) (*ChatEvent, error) {
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
		OccurredAt: time.Now().Unix(),
		Payload:    jsonChat,
	}, nil
}

func NewUpdatedEvent(c *Chat) (*ChatEvent, error) {
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
		OccurredAt: time.Now().Unix(),
		Payload:    jsonChat,
	}, nil
}

func NewDeletedEvent(roomId uuid.UUID, chatId uuid.UUID) (*ChatEvent, error) {
	p := ChatDeletedPayload{
		ID:     chatId,
		RoomID: roomId,
	}
	jsonChat, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return &ChatEvent{
		Type:       EventTypeDeleted,
		OccurredAt: time.Now().Unix(),
		Payload:    jsonChat,
	}, nil
}
