package domain

import (
	"backend/chat/event"
	"database/sql"
	"github.com/google/uuid"
)

type ChatEventRepository interface {
	Fetch(chatId uuid.UUID, tx *sql.Tx) (event.ChatEvent, error)
	Save(event *event.ChatEvent, tx *sql.Tx) (int64, error)
}
