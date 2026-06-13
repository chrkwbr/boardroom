package main

import (
	"boardroom/chat-notification"
	"boardroom/shared/infra/pubsub"
	"chat-consumer-notifier/internal"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

func main() {
	Init()
	// Kafka → Redis read model 構築
	k := pubsub.NewKafkaReader(kafkaBrokers(), "chat-events", "redis_pubsub")
	r := notification.NewChatRedisRepository(RedisClient)
	internal.NewChatNotificationPublisher(k, r).Start()

	defer func() {
		if err := RedisClient.Close(); err != nil {
			log.Println("Error closing Redis client:", err)
		}
		k.Close()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gracefully")
}

func kafkaBrokers() []string {
	brokers := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	if brokers == "" {
		return []string{"localhost:9092"}
	}
	parts := strings.Split(brokers, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return []string{"localhost:9092"}
	}
	return out
}

func redisAddr() string {
	addr := strings.TrimSpace(os.Getenv("REDIS_ADDR"))
	if addr == "" {
		return "localhost:6379"
	}
	return addr
}

func Init() {
	log.Println("Initializing consumer kafka chat")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr(),
		Password: "",
		DB:       0,
	})
}
