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
	RoomId  uuid.UUID      `json:"roomId"`
	ChatId  uuid.UUID      `json:"chatId"`
	Payload *ChatReadModel `json:"payload,omitempty"`
}

func NewChatCreatedEvent(r ChatReadModel) *ChatRedisEvent {
	return &ChatRedisEvent{
		Type:    EventTypeCreated,
		RoomId:  r.RoomID,
		ChatId:  r.ID,
		Payload: &r,
	}
}

func NewChatEditedEvent(r ChatReadModel) *ChatRedisEvent {
	return &ChatRedisEvent{
		Type:    EventTypeUpdated,
		RoomId:  r.RoomID,
		ChatId:  r.ID,
		Payload: &r,
	}
}

func NewChatDeletedEvent(r uuid.UUID, c uuid.UUID) *ChatRedisEvent {
	return &ChatRedisEvent{
		Type:   EventTypeDeleted,
		RoomId: r,
		ChatId: c,
	}
}
