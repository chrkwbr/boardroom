package main

import (
	"backend/chat/cmd/chat/api/query/internal/api"
	"backend/chat/cmd/chat/api/query/internal/repository"
	"backend/chat/cmd/chat/api/query/internal/service"
	"backend/chat/pkg/shared/readmodel"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	scylla, err := readmodel.NewChatScyllaRepository("localhost")
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
