package query

import (
	"backend/chat/command/domain"
	"backend/chat/event"
	"encoding/json"
)

type ChatResponse struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Version   int64  `json:"version"`
	Timestamp int64  `json:"date"`
}

type ChatReadModel struct {
	ID        string
	Sender    string
	Room      string
	Message   string
	Version   int64
	CreatedAt int64
	UpdatedAt int64
}

func FromPayload(chatEvent *event.ChatEvent) (*ChatReadModel, error) {
	chat := &domain.Chat{}
	if err := json.Unmarshal(chatEvent.Payload, chat); err != nil {
		return nil, err
	}
	return &ChatReadModel{
		ID:        chat.ID.String(),
		Sender:    chat.Sender,
		Room:      chat.Room,
		Message:   chat.Message,
		Version:   chat.Version,
		CreatedAt: chat.Timestamp,
		UpdatedAt: chatEvent.Timestamp,
	}, nil
}
