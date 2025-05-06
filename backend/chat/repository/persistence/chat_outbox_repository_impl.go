package persistence

import (
	"backend/chat/domain"
	"database/sql"
)

type ChatOutboxRepositoryImpl struct {
}

func NewChatOutboxRepositoryImpl() *ChatOutboxRepositoryImpl {
	return &ChatOutboxRepositoryImpl{}
}

func (impl *ChatOutboxRepositoryImpl) Save(event *domain.ChatEventOutbox, tx *sql.Tx) (int64, error) {
	insertQuery := "INSERT INTO chat_outbox (event_id, payload, created_at) VALUES ($1, $2, $3) RETURNING id"
	var id int64
	if err := tx.QueryRow(insertQuery, event.EventId, event.Payload, event.Timestamp).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
