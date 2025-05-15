package processor

import (
	"backend/chat/command/domain"
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
	chat := &domain.Chat{}
	if err := json.Unmarshal(msg, chat); err != nil {
		log.Println("Failed to unmarshal chat:", err)
		return
	}
	chatID := chat.ID
	key := fmt.Sprintf("chat:%s", chatID)
	err := p.redisClient.Set(ctx, key, msg, time.Hour*24*10).Err()
	if err != nil {
		log.Println("Error publishing to Redis:", err)
		return
	}

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
