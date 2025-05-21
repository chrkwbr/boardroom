package event

import "github.com/google/uuid"

type ChatEvent struct {
	ChatId    uuid.UUID
	EventType string
	Version   int64
	Payload   []byte
	Timestamp int64
}

const (
	ChatCreatedEvent = "chat_created"
	ChatEditedEvent  = "chat_edited"
	ChatDeletedEvent = "chat_deleted"
)
