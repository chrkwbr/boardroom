package main

import (
	apichat "backend/chat/api"
	"backend/chat/processor"
	"backend/chat/repository/persistence"
	"backend/chat/usecase"
	wschat "backend/chat/ws"
	tx "backend/infra/db"
	"backend/infra/pubsub/kafka"
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"time"
)

func main() {
	r := gin.New()
	chatDb, err := sql.Open("postgres", "host=localhost port=5432 user=boardroom password=boardroom dbname=boardroom search_path=chat sslmode=disable")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 許可するオリジン
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	pub := kafka.NewKafkaWriter([]string{"localhost:9092"})
	chatOutboxProcessor := processor.NewOutboxProcessor(
		chatDb,
		persistence.NewChatOutboxRepositoryImpl(),
		pub,
		5*time.Second)
	chatOutboxProcessor.Start()

	func() {
		if err != nil {
			panic(err)
		}
		controller := apichat.NewChatController(
			usecase.NewChatUseCase(
				persistence.NewChatRepositoryImpl(),
				persistence.NewChatOutboxRepositoryImpl(),
				tx.NewTransactionManager(chatDb),
			),
		)

		api := r.Group("/api")
		controller.RegisterRoutes(api)
	}()

	//apinotification.RegisterRoutes(api)

	ws := r.Group("/ws")
	wschat.RegisterRoutes(ws)

	r.Run()

	defer func() {
		if err := chatDb.Close(); err != nil {
			panic(err)
		}
		chatOutboxProcessor.Stop()
	}()
}
