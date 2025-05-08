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

func NewRedisProcessor(rdb *redis.Client) *RedisProcessor {
	kafkaHub, err := hub.GetHubFactory().GetHub(hub.ChatEventKafkaWs)
	if err != nil {
		log.Panicln("Failed to get hub:", err)
	}
	client := kafkaHub.CreateAndRegisterClient(10)
	p := &RedisProcessor{
		hub:         kafkaHub,
		hubClient:   client,
		redisClient: rdb,
	}

	go client.Receive(func(bytes []byte) {
		p.Process(bytes, context.Background())
	})
	return p
}

func (p *RedisProcessor) Process(msg []byte, ctx context.Context) {
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
}

func (p *RedisProcessor) Close() {
	if p.hubClient != nil {
		p.hub.UnregisterClient(p.hubClient)
		p.hubClient = nil
	}

}
