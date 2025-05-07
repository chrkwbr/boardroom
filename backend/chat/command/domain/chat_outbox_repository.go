package domain

import "database/sql"

type ChatOutboxRepository interface {
	Save(outbox *ChatEventOutbox, tx *sql.Tx) (int64, error)
	FetchUnprocessed(limit int, tx *sql.Tx) ([]*ChatEventOutbox, error)
	Delete(entity *ChatEventOutbox, tx *sql.Tx) error
}
