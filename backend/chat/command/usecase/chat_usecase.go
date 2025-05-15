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
}

func NewChatUseCase(
	chatRepository domain.ChatRepository,
	chatOutboxRepository domain.ChatOutboxRepository,
	txManager db.Transaction,
) *ChatUseCase {
	return &ChatUseCase{
		chatRepository:       chatRepository,
		chatOutboxRepository: chatOutboxRepository,
		txManager:            txManager,
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
		outbox := domain.AsOutbox(eventId, event)
		_, err = uc.chatOutboxRepository.Save(&outbox, tx)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	marshal, _ := json.Marshal(&event)

	outboxHub, err := hub.GetHubFactory().GetHub(hub.ChatEventOutbox)
	if err != nil {
		return err
	}
	outboxHub.BroadcastMessage(marshal)

	return nil
}
