package main

import (
	chatqueryApi "backend/internal/chat/query/api"
	"backend/internal/chat/query/repository"
	"backend/internal/chat/query/service"
	"backend/internal/chat/readmodel"
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
	chatQueryApi := chatqueryApi.NewChatQueryController(
		service.NewChatService(chatReadModelRepository),
	)

	api := r.Group("/api")
	chatQueryApi.RegisterRoutes(api)

	port := os.Getenv("PORT")
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

