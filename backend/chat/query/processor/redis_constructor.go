package processor

import (
	"backend/chat/command/domain"
	"backend/chat/event"
	"backend/infra/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisConstructor struct {
	subscriber  pubsub.EventSubscriber
	redisClient *redis.Client
}

func NewRedisConstructor(sub pubsub.EventSubscriber, rdb *redis.Client) *RedisConstructor {
	return &RedisConstructor{
		subscriber:  sub,
		redisClient: rdb,
	}
}

func (p *RedisConstructor) Start() {
	go func() {
		if err := p.subscriber.Subscribe("_", func(key string, value []byte) error {
			p.process(value, context.Background())
			return nil
		}); err != nil {
			log.Panicln(

				"Failed to subscribe to event:", err)
		}
		log.Println("Event subscriber started")
	}()
}

func (p *RedisConstructor) process(msg []byte, ctx context.Context) {
	chatEvent := &event.ChatEvent{}
	if err := json.Unmarshal(msg, chatEvent); err != nil {
		log.Println("Failed to unmarshal chat:", err)
		return
	}

	switch chatEvent.EventType {
	case event.ChatCreatedEvent:
		p.createReadModel(chatEvent, ctx)
	case event.ChatEditedEvent:
		p.updateReadModel(chatEvent, ctx)
	case event.ChatDeletedEvent:
		p.DeleteReadModel(chatEvent, ctx)
	}

}

func (p *RedisConstructor) createReadModel(chatEvent *event.ChatEvent, ctx context.Context) {
	key := fmt.Sprintf("chat:%s", chatEvent.ChatId)
	err := p.redisClient.Set(ctx, key, chatEvent.Payload, time.Hour*24*10).Err()
	if err != nil {
		log.Println("Error publishing to Redis:", err)
		return
	}

	chat := &domain.Chat{}
	if err := json.Unmarshal(chatEvent.Payload, chat); err != nil {
		log.Println("Failed to unmarshal chat:", err)
		return
	}
	chatID := chatEvent.ChatId

	chatRoomKey := fmt.Sprintf("chats:%v", chat.Room)
	z := redis.Z{
		Score:  float64(chat.Timestamp),
		Member: fmt.Sprintf("%s", chatID),
	}
	err = p.redisClient.ZAddNX(ctx, chatRoomKey, z).Err()
	if err != nil {
		log.Println("Error adding to Redis sorted set:", err)
	}
}

func (p *RedisConstructor) updateReadModel(chatEvent *event.ChatEvent, ctx context.Context) {
	key := fmt.Sprintf("chat:%s", chatEvent.ChatId)
	previewChat, err := p.redisClient.SetArgs(ctx, key, chatEvent.Payload, redis.SetArgs{
		Get: true,
		TTL: time.Hour * 24 * 10,
	}).Result()
	if err != nil {
		log.Println("Error updating Redis:", err)
		return
	}
	chatHistoryKey := fmt.Sprintf("chats:%v:history", chatEvent.ChatId)
	p.redisClient.LPush(ctx, chatHistoryKey, previewChat)
}

func (p *RedisConstructor) DeleteReadModel(chatEvent *event.ChatEvent, ctx context.Context) {
	chat := &domain.Chat{}
	if err := json.Unmarshal(chatEvent.Payload, chat); err != nil {
		log.Println("Failed to unmarshal chat:", err)
		return
	}
	chatRoomKey := fmt.Sprintf("chats:%v", chat.Room)
	if err := p.redisClient.ZRem(ctx, chatRoomKey, fmt.Sprintf("%s", chatEvent.ChatId)).Err(); err != nil {
		log.Println("Error removing from Redis sorted set:", err)
	}

	key := fmt.Sprintf("chat:%s", chatEvent.ChatId)
	if err := p.redisClient.Del(ctx, key).Err(); err != nil {
		log.Println("Error deleting from Redis:", err)
		return
	}
}
