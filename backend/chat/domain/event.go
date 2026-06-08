package domain

import "github.com/google/uuid"

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
	Message string    `json:"message"`
}
