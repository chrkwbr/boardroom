package main

import (
	"chat-ws/internal/handler"
	"chat-ws/internal/processor"
	"log"
	"os"
	"strings"

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
	chatWs := handler.NewChatWebSocket()
	chatWs.RegisterRoutes(ws)

	// Redis pub/sub → WebSocket への push
	processor.NewChatNotificationSubscriber(RedisClient).Start()

	defer func() {
		if err := RedisClient.Close(); err != nil {
			log.Println("Error closing Redis client:", err)
		}
	}()

	port := os.Getenv("PORT")
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func redisAddr() string {
	addr := strings.TrimSpace(os.Getenv("REDIS_ADDR"))
	if addr == "" {
		return "localhost:6379"
	}
	return addr
}

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr(),
		Password: "",
		DB:       0,
	})
}
