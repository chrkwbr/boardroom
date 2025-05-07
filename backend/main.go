package main

import (
	chatCommandApi "backend/chat/command/api"
	processor2 "backend/chat/command/processor"
	persistence2 "backend/chat/command/repository/persistence"
	"backend/chat/command/usecase"
	chatqueryApi "backend/chat/query/api"
	"backend/chat/query/processor"
	"backend/chat/query/service"
	wschat "backend/chat/ws"
	tx "backend/infra/db"
	"backend/infra/hub"
	"backend/infra/pubsub/kafka"
	"context"
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	RedisClient *redis.Client
	ChatDB      *sql.DB
)

func main() {
	Init()
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 許可するオリジン
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// == Chat API ==
	pub := kafka.NewKafkaWriter([]string{"localhost:9092"})
	chatOutboxHub := hub.NewHub()
	go chatOutboxHub.Run()
	chatOutboxProcessor := processor2.NewOutboxProcessor(
		ChatDB,
		persistence2.NewChatOutboxRepositoryImpl(),
		pub,
		chatOutboxHub)

	chatCommandApi := chatCommandApi.NewChatCommandController(
		usecase.NewChatUseCase(
			persistence2.NewChatRepositoryImpl(),
			persistence2.NewChatOutboxRepositoryImpl(),
			tx.NewTransactionManager(ChatDB),
			chatOutboxHub,
		),
	)
	chatQueryApi := chatqueryApi.NewChatQueryController(
		service.NewChatService(RedisClient),
	)

	api := r.Group("/api")
	chatCommandApi.RegisterRoutes(api)
	chatQueryApi.RegisterRoutes(api)

	//apinotification.RegisterRoutes(api)

	// == Chat WebSocket ==
	hub := hub.NewHub()
	go hub.Run()
	sub := kafka.NewKafkaReader([]string{"localhost:9092"}, "chat_messages")
	go func() {
		if err := sub.Subscribe("_", func(key string, value []byte) error {
			hub.BroadcastMessage(value)
			return nil
		}); err != nil {
			panic(err)
		}
	}()

	ws := r.Group("/ws")
	chatWs := wschat.NewChatWebSocket(hub)
	chatWs.RegisterRoutes(ws)

	rp := processor.NewRedisProcessor(hub, RedisClient)
	rp.Process(context.Background())

	defer func() {
		if err := ChatDB.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
		chatOutboxProcessor.Close()
		rp.Close()

		if err := RedisClient.Close(); err != nil {
			log.Println("Error closing Redis client:", err)
		}

	}()

	r.Run()
}

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	cdb, err := sql.Open("postgres", "host=localhost port=5432 user=boardroom password=boardroom dbname=boardroom search_path=chat sslmode=disable")
	if err != nil {
		panic(err)
	}
	ChatDB = cdb
}
