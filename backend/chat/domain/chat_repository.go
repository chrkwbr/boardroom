package domain

import "database/sql"

type ChatRepository interface {
	Save(event *ChatEvent, tx *sql.Tx) (int64, error)
}
