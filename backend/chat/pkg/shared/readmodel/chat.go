package readmodel

import "github.com/google/uuid"

type ChatReadModel struct {
	RoomID    uuid.UUID `json:"roomId"`
	ID        uuid.UUID `json:"id"`
	Message   string    `json:"message"`
	SenderID  uuid.UUID `json:"sender"`
	Version   int64     `json:"version"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt int64     `json:"updatedAt"`
}

func (c *ChatReadModel) NewUpdate(message string, occurredAt int64) *ChatReadModel {
	return &ChatReadModel{
		RoomID:    c.RoomID,
		ID:        c.ID,
		Message:   message,
		SenderID:  c.SenderID,
		Version:   c.Version + 1,
		CreatedAt: c.CreatedAt,
		UpdatedAt: occurredAt,
	}
}
