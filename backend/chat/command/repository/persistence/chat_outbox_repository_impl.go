package persistence

import (
	"backend/chat/command/domain"
	"database/sql"
)

type ChatOutboxRepositoryImpl struct {
}

func NewChatOutboxRepositoryImpl() *ChatOutboxRepositoryImpl {
	return &ChatOutboxRepositoryImpl{}
}

func (impl *ChatOutboxRepositoryImpl) Save(event *domain.ChatEventOutbox, tx *sql.Tx) (int64, error) {
	insertQuery := "INSERT INTO chat_outbox (event_id, event_type, payload, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	var id int64
	if err := tx.QueryRow(insertQuery, event.EventId, event.EventType, event.Payload, event.Timestamp).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (impl *ChatOutboxRepositoryImpl) FetchUnprocessed(limit int, tx *sql.Tx) ([]*domain.ChatEventOutbox, error) {
	query := "SELECT id, event_id, event_type, payload, created_at FROM chat_outbox LIMIT $1 FOR UPDATE SKIP LOCKED"
	rows, err := tx.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var outboxes []*domain.ChatEventOutbox
	for rows.Next() {
		var id int64
		var outbox domain.ChatEventOutbox
		if err := rows.Scan(&id, &outbox.EventId, &outbox.EventType, &outbox.Payload, &outbox.Timestamp); err != nil {
			return nil, err
		}
		outboxes = append(outboxes, &outbox)
	}

	return outboxes, nil
}

func (impl *ChatOutboxRepositoryImpl) Delete(entity *domain.ChatEventOutbox, tx *sql.Tx) error {
	deleteQuery := "DELETE FROM chat_outbox WHERE event_id = $1"
	_, err := tx.Exec(deleteQuery, entity.EventId)
	if err != nil {
		return err
	}
	return nil
}
