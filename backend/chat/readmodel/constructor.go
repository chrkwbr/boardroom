package readmodel

import (
	"backend/chat/domain"
	"backend/infra/pubsub"
	"context"
	"encoding/json"
	"log"
)

type ChatHistory struct {
	ChatId    string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

type RedisConstructor struct {
	subscriber pubsub.EventSubscriber
	repo       *ChatRedisRepository
}

func NewRedisConstructor(sub pubsub.EventSubscriber, repo *ChatRedisRepository) *RedisConstructor {
	return &RedisConstructor{
		subscriber: sub,
		repo:       repo,
	}
}

func (rc *RedisConstructor) Start() {
	go func() {
		if err := rc.subscriber.Subscribe("_", func(key string, value []byte) error {
			rc.process(value, context.Background())
			return nil
		}); err != nil {
			log.Panicln(

				"Failed to subscribe to event:", err)
		}
		log.Println("Event subscriber started")
	}()
}

func (rc *RedisConstructor) process(msg []byte, ctx context.Context) {
	chatEvent := &domain.ChatEvent{}
	if err := json.Unmarshal(msg, chatEvent); err != nil {
		log.Println("Failed to unmarshal chat:", err)
		return
	}

	switch chatEvent.Type {
	case domain.EventTypeCreated:
		rc.createReadModel(ctx, chatEvent)
	case domain.EventTypeUpdated:
		rc.updateReadModel(ctx, chatEvent)
	case domain.EventTypeDeleted:
		rc.DeleteReadModel(ctx, chatEvent)
	}

}

func (rc *RedisConstructor) createReadModel(ctx context.Context, event *domain.ChatEvent) {
	chat := &domain.ChatCreatedPayload{}
	if err := json.Unmarshal(event.Payload, chat); err != nil {
		log.Println("Failed to unmarshal chat:", err)
		return
	}
	s := User{
		ID:   chat.SenderID,
		Name: "test name",
		Icon: "test icon",
	}

	r := &ChatReadModel{
		ID:        chat.ID.String(),
		Sender:    s,
		RoomID:    chat.RoomID.String(),
		Message:   chat.Message,
		Version:   chat.Version,
		CreatedAt: event.OccurredAt,
		UpdatedAt: event.OccurredAt,
	}

	if err := rc.repo.SetChat(ctx, r); err != nil {
		log.Println("Failed to save chat read model:", err)
		return
	}

	if err := rc.repo.LPushHistory(ctx, r); err != nil {
		log.Println("Failed to push chat history:", err)
		return
	}

	if err := rc.repo.ZAddNXRoomChatIds(ctx, r); err != nil {
		log.Println("Failed to save room chat IDs:", err)
		return
	}
}

func (rc *RedisConstructor) updateReadModel(ctx context.Context, chatEvent *domain.ChatEvent) {
	p := &domain.ChatEditedPayload{}
	if err := json.Unmarshal(chatEvent.Payload, p); err != nil {
		log.Println("Failed to unmarshal chat:", err)
		return
	}

	orig, err := rc.repo.GetChat(ctx, p.ID)
	if err != nil {
		log.Println("Failed to get chat:", err)
		return
	}

	edited, err := orig.NewUpdate(p.Message, chatEvent.OccurredAt)

	if err := rc.repo.SetChat(ctx, edited); err != nil {
		log.Println("Failed to update chat read model:", err)
		return
	}
	if err := rc.repo.LPushHistory(ctx, edited); err != nil {
		log.Println("Failed to push chat history:", err)
		return
	}
}

func (rc *RedisConstructor) DeleteReadModel(ctx context.Context, chatEvent *domain.ChatEvent) {
	p := &domain.ChatDeletedPayload{}
	if err := json.Unmarshal(chatEvent.Payload, p); err != nil {
		log.Println("Failed to unmarshal chat:", err)
		return
	}

	orig, err := rc.repo.GetChat(ctx, p.ID)
	if err != nil {
		log.Println("Failed to get chat:", err)
		return
	}

	if err := rc.repo.ZRemRoomChatIds(ctx, orig.RoomID, p.ID); err != nil {
		log.Println("Failed to remove chat ID from room:", err)
		return
	}

	if err := rc.repo.DelChat(ctx, p.ID); err != nil {
		log.Println("Failed to delete chat read model:", err)
		return
	}

	if err := rc.repo.DelHistory(ctx, p.ID); err != nil {
		log.Println("Failed to delete chat history:", err)
		return
	}
}
