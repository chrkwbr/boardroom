package processor

import (
	"backend/chat/event"
	"backend/chat/query"
	"backend/chat/query/repository"
	"backend/infra/pubsub"
	"context"
	"encoding/json"
	"log"
)

type RedisConstructor struct {
	subscriber pubsub.EventSubscriber
	repository *repository.ChatReadModelRepository
}

func NewRedisConstructor(sub pubsub.EventSubscriber, modelRepository *repository.ChatReadModelRepository) *RedisConstructor {
	return &RedisConstructor{
		subscriber: sub,
		repository: modelRepository,
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
	chatEvent := &event.ChatEvent{}
	if err := json.Unmarshal(msg, chatEvent); err != nil {
		log.Println("Failed to unmarshal chat:", err)
		return
	}

	switch chatEvent.EventType {
	case event.ChatCreatedEvent:
		rc.createReadModel(ctx, chatEvent)
	case event.ChatEditedEvent:
		rc.updateReadModel(ctx, chatEvent)
	case event.ChatDeletedEvent:
		rc.DeleteReadModel(ctx, chatEvent)
	}

}

func (rc *RedisConstructor) createReadModel(ctx context.Context, chatEvent *event.ChatEvent) {
	readModel, err := rc.repository.SetChat(ctx, chatEvent)
	if err != nil {
		log.Println("Failed to save chat read model:", err)
		return
	}

	if err := rc.repository.LPushHistory(ctx, readModel); err != nil {
		log.Println("Failed to push chat history:", err)
		return
	}

	if err := rc.repository.ZAddNXRoomChatIds(ctx, readModel); err != nil {
		log.Println("Failed to save room chat IDs:", err)
		return
	}
}

func (rc *RedisConstructor) updateReadModel(ctx context.Context, chatEvent *event.ChatEvent) {
	readModel, err := rc.repository.SetChat(ctx, chatEvent)
	if err != nil {
		log.Println("Failed to update chat read model:", err)
		return
	}
	if err := rc.repository.LPushHistory(ctx, readModel); err != nil {
		log.Println("Failed to push chat history:", err)
		return
	}
}

func (rc *RedisConstructor) DeleteReadModel(ctx context.Context, chatEvent *event.ChatEvent) {
	readModel, err := query.FromPayload(chatEvent)
	if err != nil {
		log.Println("Failed to convert payload to read model:", err)
		return
	}
	if err := rc.repository.ZRemRoomChatIds(ctx, readModel.Room, chatEvent.ChatId); err != nil {
		log.Println("Failed to remove chat ID from room:", err)
		return
	}

	if err := rc.repository.DelChat(ctx, chatEvent.ChatId); err != nil {
		log.Println("Failed to delete chat read model:", err)
		return
	}

	if err := rc.repository.DelHistory(ctx, chatEvent.ChatId); err != nil {
		log.Println("Failed to delete chat history:", err)
		return
	}
}
