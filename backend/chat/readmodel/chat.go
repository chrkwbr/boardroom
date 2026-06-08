package readmodel

type ChatReadModel struct {
	ID        string
	Sender    User
	RoomID    string
	Message   string
	Version   int64
	CreatedAt int64
	UpdatedAt int64
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
