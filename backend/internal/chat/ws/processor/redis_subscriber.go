package processor

import (
	"backend/internal/shared/infra/hub"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type ChatRedisSubscriber struct {
	redis *redis.Client
}

func NewChatRedisSubscriber(redis *redis.Client) *ChatRedisSubscriber {
	return &ChatRedisSubscriber{
		redis: redis,
	}
}

func (k *ChatRedisSubscriber) Start() {
	h, err := hub.GetHubFactory().GetHub(hub.ChatEventWsPusher)
	if err != nil {
		log.Panicln("Failed to get hub:", err)
	}
	go func() {
		pubsub := k.redis.PSubscribe(context.Background(), "chat:room:*:updates")
		ch := pubsub.Channel()
		for msg := range ch {
			h.BroadcastMessage([]byte(msg.Payload))
		}
	}()

}
