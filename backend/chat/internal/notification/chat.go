package notification

import "github.com/google/uuid"

type Chat struct {
	RoomID    uuid.UUID `json:"roomId"`
	ID        uuid.UUID `json:"id"`
	Message   string    `json:"message"`
	Sender    User      `json:"sender"`
	Version   int64     `json:"version"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt int64     `json:"updatedAt"`
}

func (c *Chat) NewUpdate(message string, occurredAt int64) *Chat {
	return &Chat{
		RoomID:    c.RoomID,
		ID:        c.ID,
		Message:   message,
		Sender:    c.Sender,
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
		Sender:    c.Sender,
		Version:   c.Version + 1,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
