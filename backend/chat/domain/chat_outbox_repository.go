package domain

import "database/sql"

type ChatOutboxRepository interface {
	Save(outbox *ChatEventOutbox, tx *sql.Tx) (int64, error)
}
