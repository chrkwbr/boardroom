package main

import (
	"boardroom/shared/infra/pubsub/kafka"
	"chat-api-command/internal/api"
	"chat-api-command/internal/usecase"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// == Chat API ==
	eventPublisher := kafka.NewKafkaWriter([]string{"localhost:9092"})
	chatCommandApi := api.NewChatCommandController(
		usecase.NewChatUseCase(
			eventPublisher,
		),
	)

	api := r.Group("/api")
	chatCommandApi.RegisterRoutes(api)

	defer func() {
		eventPublisher.Close()
	}()

	port := os.Getenv("PORT")
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
