package usecase

import (
	domain2 "backend/chat/domain"
	"backend/infra/db"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type ChatUseCase struct {
	chatRepository       domain2.ChatRepository
	chatOutboxRepository domain2.ChatOutboxRepository
	txManager            db.Transaction
}

func NewChatUseCase(
	chatRepository domain2.ChatRepository,
	chatOutboxRepository domain2.ChatOutboxRepository,
	txManager db.Transaction,
) *ChatUseCase {
	return &ChatUseCase{
		chatRepository:       chatRepository,
		chatOutboxRepository: chatOutboxRepository,
		txManager:            txManager,
	}
}

func (uc *ChatUseCase) CreateChat(sender string, room string, message string) error {
	chat := &domain2.Chat{
		ID:        uuid.New(),
		Sender:    sender,
		Room:      room,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	return uc.txManager.RunWithTx(func(tx *sql.Tx) error {
		event := chat.AsCreateEvent()
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
	})
}
