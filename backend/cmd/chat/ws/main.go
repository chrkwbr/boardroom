package main

import (
	"backend/internal/chat/readmodel"
	wschat "backend/internal/chat/ws/handler"
	"backend/internal/chat/ws/processor"
	"backend/internal/shared/infra/pubsub/kafka"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

func main() {
	Init()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	ws := r.Group("/ws")
	chatWs := wschat.NewChatWebSocket()
	chatWs.RegisterRoutes(ws)

	// Redis pub/sub → WebSocket への push
	processor.NewChatRedisSubscriber(RedisClient).Start()

	// Kafka → Redis read model 構築
	subscriberRedisConstructor := kafka.NewKafkaReader([]string{"localhost:9092"}, "chat-events", "redis_constructor")
	readmodelRedis := readmodel.NewChatRedisRepository(RedisClient)
	readmodel.NewRedisConstructor(subscriberRedisConstructor, readmodelRedis).Start()

	defer func() {
		if err := RedisClient.Close(); err != nil {
			log.Println("Error closing Redis client:", err)
		}
		subscriberRedisConstructor.Close()
	}()

	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
