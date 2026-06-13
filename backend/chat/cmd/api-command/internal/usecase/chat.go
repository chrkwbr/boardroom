package usecase

import (
	"boardroom/chat-shared/domain"
	"boardroom/chat-shared/infra/pubsub"
	"encoding/json"

	"github.com/google/uuid"
)

type ChatUseCase struct {
	publisher pubsub.EventPublisher
}

func NewChatUseCase(
	publisher pubsub.EventPublisher,
) *ChatUseCase {
	return &ChatUseCase{
		publisher: publisher,
	}
}

func (uc *ChatUseCase) CreateChat(senderID uuid.UUID, roomID uuid.UUID, message string) error {
	// ToDo validation
	// sender が room に属しているか
	chat := domain.NewChat(senderID, roomID, message)
	event, err := domain.NewCreatedEvent(chat)
	if err != nil {
		return err
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	// 同じキーを持つイベントは同じパーティションに順番保証されるので、キーは Room
	if err := uc.publisher.Publish("chat-events", roomID.String(), eventJSON); err != nil {
		return err
	}
	return nil
}

func (uc *ChatUseCase) EditChat(roomID uuid.UUID, chatId uuid.UUID, message string) error {
	// ToDo validation
	// sender が room に属しているか
	// chatId の投稿者が同一か
	chat := domain.NewEditedChat(chatId, uuid.Nil, roomID, message)
	event, err := domain.NewUpdatedEvent(chat)
	if err != nil {
		return err
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := uc.publisher.Publish("chat-events", roomID.String(), eventJSON); err != nil {
	}
	return nil
}

func (uc *ChatUseCase) DeleteChat(roomID uuid.UUID, chatId uuid.UUID) error {
	// ToDo validation
	// sender が room に属しているか
	// chatId の投稿者が同一か
	del, err := domain.NewDeletedEvent(roomID, chatId)
	if err != nil {
		return err
	}

	eventJSON, err := json.Marshal(del)
	if err != nil {
		return err
	}
	if err := uc.publisher.Publish("chat-events", roomID.String(), eventJSON); err != nil {
		return err
	}
	return nil
}
