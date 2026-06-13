package usecase

import (
	"backend/chat/pkg/shared/domain"
	"backend/chat/pkg/shared/infra/pubsub"
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
	event := chat.NewCreatedEvent()
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
	event := chat.NewUpdatedEvent()
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := uc.publisher.Publish("chat-events", roomID.String(), eventJSON); err != nil {
	}
	return nil
}

func (uc *ChatUseCase) DeleteChat(chatId uuid.UUID) interface{} {
	//var marshaledEvent []byte
	//if err := uc.txManager.RunWithTx(func(tx *sql.Tx) error {
	//	chatEvent, err := uc.chatRepository.Fetch(chatId, tx)
	//	if err != nil {
	//		return err
	//	}
	//	var chat domain.Chat
	//	if err := json.Unmarshal(chatEvent.Payload, &chat); err != nil {
	//		return err
	//	}
	//	chat.Version = chat.Version + 1
	//
	//	event := chat.AsDeleteEvent()
	//	eventId, err := uc.chatRepository.Save(&event, tx)
	//	if err != nil {
	//		return err
	//	}
	//	outbox := domain.AsOutbox(eventId, event)
	//	_, err = uc.chatOutboxRepository.Save(&outbox, tx)
	//	if err != nil {
	//		return err
	//	}
	//	marshaledEvent, _ = json.Marshal(&event)
	//	return nil
	//}); err != nil {
	//	return err
	//}
	//outboxHub, err := hub.GetHubFactory().GetHub(hub.ChatEventOutbox)
	//if err != nil {
	//	return err
	//}
	//outboxHub.BroadcastMessage(marshaledEvent)
	return nil
}
