package main

import (
	chatqueryApi "backend/internal/chat/query/api"
	"backend/internal/chat/query/repository"
	"backend/internal/chat/query/service"
	"log"
	"os"

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

	// == Chat API ==
	chatReadModelRepository := repository.NewChatReadModelRepository(RedisClient)
	chatQueryApi := chatqueryApi.NewChatQueryController(
		service.NewChatService(chatReadModelRepository),
	)

	api := r.Group("/api")
	chatQueryApi.RegisterRoutes(api)

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

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
