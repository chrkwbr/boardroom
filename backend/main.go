package main

import (
	apichat "backend/chat/api"
	"backend/chat/processor"
	"backend/chat/repository/persistence"
	"backend/chat/usecase"
	wschat "backend/chat/ws"
	tx "backend/infra/db"
	"backend/infra/hub"
	"backend/infra/pubsub/kafka"
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.New()
	chatDb, err := sql.Open("postgres", "host=localhost port=5432 user=boardroom password=boardroom dbname=boardroom search_path=chat sslmode=disable")
	if err != nil {
		panic(err)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 許可するオリジン
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	pub := kafka.NewKafkaWriter([]string{"localhost:9092"})
	chatOutboxHub := hub.NewHub()
	go chatOutboxHub.Run()
	chatOutboxProcessor := processor.NewOutboxProcessor(
		chatDb,
		persistence.NewChatOutboxRepositoryImpl(),
		pub,
		chatOutboxHub)

	func() {
		controller := apichat.NewChatController(
			usecase.NewChatUseCase(
				persistence.NewChatRepositoryImpl(),
				persistence.NewChatOutboxRepositoryImpl(),
				tx.NewTransactionManager(chatDb),
				chatOutboxHub,
			),
		)

		api := r.Group("/api")
		controller.RegisterRoutes(api)

	}()

	//apinotification.RegisterRoutes(api)

	ws := r.Group("/ws")
	wschat.RegisterRoutes(ws)

	defer func() {
		if err := chatDb.Close(); err != nil {
			panic(err)
		}
		chatOutboxProcessor.Close()
	}()

	r.Run()
}
