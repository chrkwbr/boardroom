package processor

import (
	"backend/chat/command/domain"
	"backend/infra/hub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisProcessor struct {
	hub         *hub.Hub
	hubClient   *hub.Client
	redisClient *redis.Client
}

func NewRedisProcessor(h *hub.Hub, rdb *redis.Client) *RedisProcessor {
	return &RedisProcessor{
		hub:         h,
		hubClient:   hub.NewClient(10),
		redisClient: rdb,
	}
}

func (p *RedisProcessor) Process(ctx context.Context) {
	p.hub.RegisterClient(p.hubClient)
	go p.hubClient.Receive(func(msg []byte) {
		chat := &domain.Chat{}
		if err := json.Unmarshal(msg, chat); err != nil {
			log.Println("Failed to unmarshal chat:", err)
			return
		}
		key := fmt.Sprintf("chat:room:%s", chat.Room)
		// ToDo check if the key exists
		err := p.redisClient.RPush(ctx, key, msg).Err()
		if err != nil {
			log.Println("Error publishing to Redis:", err)
			return
		}
	})
}

func (p *RedisProcessor) Close() {
	if p.hubClient != nil {
		p.hub.UnregisterClient(p.hubClient)
		p.hubClient = nil
	}

}
