package usecase

import (
	"backend/internal/chat/domain"
	"backend/internal/shared/infra/pubsub"
	"encoding/json"

	"github.com/google/uuid"
)

type ChatUseCase struct {
	//chatRepository       domain.ChatEventRepository
	//chatOutboxRepository domain.ChatOutboxRepository
	//txManager            db.Transaction
	publisher pubsub.EventPublisher
}

func NewChatUseCase(
	//chatRepository domain.ChatEventRepository,
	//chatOutboxRepository domain.ChatOutboxRepository,
	//txManager db.Transaction,
	publisher pubsub.EventPublisher,
) *ChatUseCase {
	return &ChatUseCase{
		//chatRepository:       chatRepository,
		//chatOutboxRepository: chatOutboxRepository,
		//txManager:            txManager,
		publisher: publisher,
	}
}

func (uc *ChatUseCase) CreateChat(senderID uuid.UUID, roomID uuid.UUID, message string) error {
	// ToDo get From Redis
	user := domain.User{
		ID:   senderID,
		Name: "Sender",
	}
	room := domain.Room{
		ID:   roomID,
		Name: "Room",
	}

	chat := domain.NewChat(user, room, message)
	event := chat.NewCreatedEvent()
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	// 同じキーを持つイベントは同じパーティションに順番保証されるので、キーは Room
	if err := uc.publisher.Publish("chat-events", roomID.String(), eventJSON); err != nil {
		return err
	}
	//if err := uc.txManager.RunWithTx(func(tx *sql.Tx) error {
	//	eventId, err := uc.chatRepository.Save(&event, tx)
	//	if err != nil {
	//		return err
	//	}
	//	outbox := domain.AsOutbox(eventId, event)
	//	_, err = uc.chatOutboxRepository.Save(&outbox, tx)
	//	if err != nil {
	//		return err
	//	}
	//	return nil
	//}); err != nil {
	//	return err
	//}
	//
	//marshal, _ := json.Marshal(&event)
	//
	//outboxHub, err := hub.GetHubFactory().GetHub(hub.ChatEventOutbox)
	//if err != nil {
	//	return err
	//}
	//outboxHub.BroadcastMessage(marshal)

	return nil
}

func (uc *ChatUseCase) EditChat(chatId uuid.UUID, message string) error {
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
	//	editedChat := chat.Edit(message)
	//
	//	event := editedChat.AsEditEvent()
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
