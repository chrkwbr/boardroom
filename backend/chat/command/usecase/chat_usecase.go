package usecase

import (
	"backend/chat/command/domain"
	"backend/infra/db"
	"backend/infra/hub"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type ChatUseCase struct {
	chatRepository       domain.ChatRepository
	chatOutboxRepository domain.ChatOutboxRepository
	txManager            db.Transaction
	hub                  *hub.Hub
}

func NewChatUseCase(
	chatRepository domain.ChatRepository,
	chatOutboxRepository domain.ChatOutboxRepository,
	txManager db.Transaction,
	hub *hub.Hub,
) *ChatUseCase {
	return &ChatUseCase{
		chatRepository:       chatRepository,
		chatOutboxRepository: chatOutboxRepository,
		txManager:            txManager,
		hub:                  hub,
	}
}

func (uc *ChatUseCase) CreateChat(sender string, room string, message string) error {
	chat := &domain.Chat{
		ID:        uuid.New(),
		Sender:    sender,
		Room:      room,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	event := chat.AsCreateEvent()
	if err := uc.txManager.RunWithTx(func(tx *sql.Tx) error {
		eventId, err := uc.chatRepository.Save(&event, tx)
		if err != nil {
			return err
		}
		outbox := event.AsOutbox(eventId)
		_, err = uc.chatOutboxRepository.Save(&outbox, tx)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	marshal, _ := json.Marshal(&event)
	uc.hub.BroadcastMessage(marshal)
	return nil
}
