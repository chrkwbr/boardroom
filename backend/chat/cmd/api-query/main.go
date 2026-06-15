package main

import (
	"boardroom/chat-readmodel"
	"chat-api-query/internal/api"
	"chat-api-query/internal/repository"
	"chat-api-query/internal/service"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func scyllaHosts() []string {
	hosts := strings.TrimSpace(os.Getenv("SCYLLA_HOST"))
	if hosts == "" {
		return []string{"localhost"}
	}
	parts := strings.Split(hosts, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return []string{"localhost"}
	}
	return out
}

func connectScyllaWithRetry() (*readmodel.ChatScyllaRepository, error) {
	const maxAttempts = 30
	const interval = 2 * time.Second

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		repo, err := readmodel.NewChatScyllaRepository(scyllaHosts()...)
		if err == nil {
			return repo, nil
		}
		lastErr = err
		log.Printf("Scylla connect attempt %d/%d failed: %v", attempt, maxAttempts, err)
		time.Sleep(interval)
	}
	return nil, lastErr
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	scylla, err := connectScyllaWithRetry()
	if err != nil {
		log.Fatal("Failed to connect to ScyllaDB:", err)
	}
	defer scylla.Close()

	chatReadModelRepository := repository.NewChatScyllaQueryRepository(scylla)
	chatQueryApi := api.NewChatQueryController(
		service.NewChatService(chatReadModelRepository),
	)

	api := r.Group("/api")
	chatQueryApi.RegisterRoutes(api)

	port := os.Getenv("PORT")
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
