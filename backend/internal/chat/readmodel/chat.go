package readmodel

import "github.com/google/uuid"

type ChatReadModel struct {
	ID        uuid.UUID `json:"id"`
	Sender    User      `json:"sender"`
	RoomID    uuid.UUID `json:"roomId"`
	Message   string    `json:"message"`
	Version   int64     `json:"version"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt int64     `json:"updatedAt"`
}

func (c *ChatReadModel) NewUpdate(message string, occurredAt int64) (*ChatReadModel, error) {
	return &ChatReadModel{
		ID:        c.ID,
		Sender:    c.Sender,
		RoomID:    c.RoomID,
		Message:   message,
		Version:   c.Version + 1,
		CreatedAt: c.CreatedAt,
		UpdatedAt: occurredAt,
	}, nil
}
