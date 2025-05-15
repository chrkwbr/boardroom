package persistence

import (
	"backend/chat/command/domain"
	"backend/chat/event"
	"database/sql"
)

type ChatRepositoryImpl struct {
}

func NewChatRepositoryImpl() domain.ChatRepository {
	return &ChatRepositoryImpl{}
}

func (impl *ChatRepositoryImpl) Save(event *event.ChatEvent, tx *sql.Tx) (int64, error) {
	// ToDo 楽観排他制御
	insertQuery := "INSERT INTO chat_events (chat_id, event_type, version, payload, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var id int64
	if err := tx.QueryRow(insertQuery, event.ChatId, event.EventType, event.Version, event.Payload, event.Timestamp).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
