package readmodel

import (
	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeCreated EventType = "created"
	EventTypeUpdated EventType = "updated"
	EventTypeDeleted EventType = "deleted"
)

type ChatRedisEvent struct {
	Type    EventType      `json:"type"`
	RoomID  uuid.UUID      `json:"roomID"`
	ChatID  uuid.UUID      `json:"chatID"`
	Payload *ChatReadModel `json:"payload,omitempty"`
}

func NewChatCreatedEvent(r ChatReadModel) *ChatRedisEvent {
	return &ChatRedisEvent{
		Type:    EventTypeCreated,
		RoomID:  r.RoomID,
		ChatID:  r.ID,
		Payload: &r,
	}
}

func NewChatEditedEvent(r ChatReadModel) *ChatRedisEvent {
	return &ChatRedisEvent{
		Type:    EventTypeUpdated,
		RoomID:  r.RoomID,
		ChatID:  r.ID,
		Payload: &r,
	}
}

func NewChatDeletedEvent(r uuid.UUID, c uuid.UUID) *ChatRedisEvent {
	return &ChatRedisEvent{
		Type:   EventTypeDeleted,
		RoomID: r,
		ChatID: c,
	}
}
