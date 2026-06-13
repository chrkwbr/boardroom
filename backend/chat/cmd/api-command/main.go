package main

import (
	"boardroom/shared/infra/pubsub"
	"chat-api-command/internal/api"
	"chat-api-command/internal/usecase"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

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

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// == Chat API ==
	eventPublisher := pubsub.NewKafkaWriter(kafkaBrokers())
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
