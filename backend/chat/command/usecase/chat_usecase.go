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
	chatRepository       domain.ChatEventRepository
	chatOutboxRepository domain.ChatOutboxRepository
	txManager            db.Transaction
}

func NewChatUseCase(
	chatRepository domain.ChatEventRepository,
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
		Version:   1,
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

func (uc *ChatUseCase) EditChat(chatId uuid.UUID, message string) error {
	var marshaledEvent []byte
	if err := uc.txManager.RunWithTx(func(tx *sql.Tx) error {
		chatEvent, err := uc.chatRepository.Fetch(chatId, tx)
		if err != nil {
			return err
		}
		var chat domain.Chat
		if err := json.Unmarshal(chatEvent.Payload, &chat); err != nil {
			return err
		}
		chat.Message = message
		chat.Version = chat.Version + 1

		event := chat.AsEditEvent()
		eventId, err := uc.chatRepository.Save(&event, tx)
		if err != nil {
			return err
		}
		outbox := domain.AsOutbox(eventId, event)
		_, err = uc.chatOutboxRepository.Save(&outbox, tx)
		if err != nil {
			return err
		}
		marshaledEvent, _ = json.Marshal(&event)
		return nil
	}); err != nil {
		return err
	}
	outboxHub, err := hub.GetHubFactory().GetHub(hub.ChatEventOutbox)
	if err != nil {
		return err
	}
	outboxHub.BroadcastMessage(marshaledEvent)

	return nil
}

func (uc *ChatUseCase) DeleteChat(chatId uuid.UUID) interface{} {
	var marshaledEvent []byte
	if err := uc.txManager.RunWithTx(func(tx *sql.Tx) error {
		chatEvent, err := uc.chatRepository.Fetch(chatId, tx)
		if err != nil {
			return err
		}
		var chat domain.Chat
		if err := json.Unmarshal(chatEvent.Payload, &chat); err != nil {
			return err
		}
		chat.Version = chat.Version + 1

		event := chat.AsDeleteEvent()
		eventId, err := uc.chatRepository.Save(&event, tx)
		if err != nil {
			return err
		}
		outbox := domain.AsOutbox(eventId, event)
		_, err = uc.chatOutboxRepository.Save(&outbox, tx)
		if err != nil {
			return err
		}
		marshaledEvent, _ = json.Marshal(&event)
		return nil
	}); err != nil {
		return err
	}
	outboxHub, err := hub.GetHubFactory().GetHub(hub.ChatEventOutbox)
	if err != nil {
		return err
	}
	outboxHub.BroadcastMessage(marshaledEvent)
	return nil
}
