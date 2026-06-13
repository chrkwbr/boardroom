package readmodel

import "github.com/google/uuid"

type Chat struct {
	RoomID    uuid.UUID `json:"roomId"`
	ID        uuid.UUID `json:"id"`
	Message   string    `json:"message"`
	SenderID  uuid.UUID `json:"sender"`
	Version   int64     `json:"version"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt int64     `json:"updatedAt"`
}

func (c *Chat) NewUpdate(message string, occurredAt int64) *Chat {
	return &Chat{
		RoomID:    c.RoomID,
		ID:        c.ID,
		Message:   message,
		SenderID:  c.SenderID,
		Version:   c.Version + 1,
		CreatedAt: c.CreatedAt,
		UpdatedAt: occurredAt,
	}
}

func (c *Chat) NewDelete() *Chat {
	return &Chat{
		RoomID:    c.RoomID,
		ID:        c.ID,
		Message:   "",
		SenderID:  c.SenderID,
		Version:   c.Version + 1,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

type Status = string

const (
	Created Status = "created"
	Edited  Status = "edited"
	Deleted Status = "deleted"
)
