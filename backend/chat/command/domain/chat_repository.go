package domain

import (
	"backend/chat/event"
	"database/sql"
)

type ChatRepository interface {
	Save(event *event.ChatEvent, tx *sql.Tx) (int64, error)
}
