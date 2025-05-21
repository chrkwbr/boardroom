package persistence

import (
	"backend/chat/command/domain"
	"backend/chat/event"
	"database/sql"
	"github.com/google/uuid"
)

type ChatEventRepositoryImpl struct {
}

func (impl *ChatEventRepositoryImpl) Fetch(chatId uuid.UUID, tx *sql.Tx) (event.ChatEvent, error) {
	query := "SELECT chat_id, event_type, version, payload, created_at  FROM chat_events WHERE chat_id = $1 AND version = (SELECT MAX(version) FROM chat_events WHERE chat_id = $1)"
	var chat event.ChatEvent
	if err := tx.QueryRow(query, chatId).Scan(
		&chat.ChatId,
		&chat.EventType,
		&chat.Version,
		&chat.Payload,
		&chat.Timestamp,
	); err != nil {
		return event.ChatEvent{}, err
	}
	return chat, nil
}

func NewChatRepositoryImpl() domain.ChatEventRepository {
	return &ChatEventRepositoryImpl{}
}

func (impl *ChatEventRepositoryImpl) Save(event *event.ChatEvent, tx *sql.Tx) (int64, error) {
	// ToDo 楽観排他制御
	insertQuery := "INSERT INTO chat_events (chat_id, event_type, version, payload, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var id int64
	if err := tx.QueryRow(insertQuery, event.ChatId, event.EventType, event.Version, event.Payload, event.Timestamp).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
