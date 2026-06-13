package processor

import (
	"boardroom/chat-shared/infra/hub"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type ChatNotificationSubscriber struct {
	redis *redis.Client
}

func NewChatNotificationSubscriber(redis *redis.Client) *ChatNotificationSubscriber {
	return &ChatNotificationSubscriber{
		redis: redis,
	}
}

func (k *ChatNotificationSubscriber) Start() {
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
