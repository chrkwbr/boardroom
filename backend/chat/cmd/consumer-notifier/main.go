package main

import (
	"boardroom/chat-shared/infra/pubsub/kafka"
	"boardroom/chat-shared/readmodel"
	"chat-consumer-notifier/internal"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

func main() {
	Init()
	// Kafka → Redis read model 構築
	k := kafka.NewKafkaReader([]string{"localhost:9092"}, "chat-events", "redis_pubsub")
	r := readmodel.NewChatRedisRepository(RedisClient)
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

func Init() {
	log.Println("Initializing consumer kafka chat")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
