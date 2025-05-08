package processor

import (
	"backend/chat/command/domain"
	"backend/infra/pubsub"
	"backend/infra/pubsub/kafka"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

type KafkaRedisProcessor struct {
	subscriber  pubsub.EventSubscriber
	redisClient *redis.Client
}

func NewRedisProcessor(rdb *redis.Client) *KafkaRedisProcessor {
	return &KafkaRedisProcessor{
		subscriber:  kafka.NewKafkaReader([]string{"localhost:9092"}, "chat_messages", "redis_constructor"),
		redisClient: rdb,
	}
}

func (p *KafkaRedisProcessor) Start() {
	go func() {
		if err := p.subscriber.Subscribe("_", func(key string, value []byte) error {
			p.process(value, context.Background())
			return nil
		}); err != nil {
			log.Panicln("Failed to subscribe to Kafka:", err)
		}
		log.Println("Kafka subscriber started")
	}()
}

func (p *KafkaRedisProcessor) process(msg []byte, ctx context.Context) {
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
