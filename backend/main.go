package main

import (
	apichat "backend/chat/api"
	"backend/chat/repository/persistence"
	"backend/chat/usecase"
	wschat "backend/chat/ws"
	tx "backend/infra/db"
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 許可するオリジン
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	func() {
		//pub := kafka.NewKafkaWriter([]string{"localhost:9092"})
		db, err := sql.Open("postgres", "host=localhost port=5432 user=boardroom password=boardroom dbname=boardroom search_path=chat sslmode=disable")
		if err != nil {
			panic(err)
		}
		controller := apichat.NewChatController(
			usecase.NewChatUseCase(
				persistence.NewChatRepositoryImpl(),
				persistence.NewChatOutboxRepositoryImpl(),
				tx.NewTransactionManager(db),
			),
		)

		api := r.Group("/api")
		controller.RegisterRoutes(api)
		//
		//h := &apichat.MessageHandler{Pub: pub}
		//h.RegisterRoutes(api)
	}()

	//apinotification.RegisterRoutes(api)

	ws := r.Group("/ws")
	wschat.RegisterRoutes(ws)

	r.Run()
}
