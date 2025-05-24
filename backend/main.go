package main

import (
	chatCommandApi "backend/chat/command/api"
	processor2 "backend/chat/command/processor"
	persistence2 "backend/chat/command/repository/persistence"
	"backend/chat/command/usecase"
	chatqueryApi "backend/chat/query/api"
	processor3 "backend/chat/query/processor"
	"backend/chat/query/repository"
	"backend/chat/query/service"
	wschat "backend/chat/ws/handler"
	"backend/chat/ws/processor"
	tx "backend/infra/db"
	"backend/infra/pubsub/kafka"
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
	eventPublisher := kafka.NewKafkaWriter([]string{"localhost:9092"})
	chatOutboxProcessor := processor2.NewOutboxProcessor(
		ChatDB,
		persistence2.NewChatOutboxRepositoryImpl(),
		eventPublisher)

	chatCommandApi := chatCommandApi.NewChatCommandController(
		usecase.NewChatUseCase(
			persistence2.NewChatRepositoryImpl(),
			persistence2.NewChatOutboxRepositoryImpl(),
			tx.NewTransactionManager(ChatDB),
		),
	)
	chatReadModelRepository := repository.NewChatReadModelRepository(RedisClient)
	chatQueryApi := chatqueryApi.NewChatQueryController(
		service.NewChatService(chatReadModelRepository),
	)

	api := r.Group("/api")
	chatCommandApi.RegisterRoutes(api)
	chatQueryApi.RegisterRoutes(api)

	//apinotification.RegisterRoutes(api)

	// == Chat WebSocket ==
	subscriberWebsocketProcessor := kafka.NewKafkaReader([]string{"localhost:9092"}, "chat_messages", "websocket_processor")
	processor.NewWsChatEventPusher(subscriberWebsocketProcessor).Start()

	ws := r.Group("/ws")
	chatWs := wschat.NewChatWebSocket()
	chatWs.RegisterRoutes(ws)

	subscriberRedisConstructor := kafka.NewKafkaReader([]string{"localhost:9092"}, "chat_messages", "redis_constructor")
	processor3.NewRedisConstructor(subscriberRedisConstructor, chatReadModelRepository).Start()

	defer func() {
		if err := ChatDB.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
		chatOutboxProcessor.Close()

		if err := RedisClient.Close(); err != nil {
			log.Println("Error closing Redis client:", err)
		}

		eventPublisher.Close()
		subscriberWebsocketProcessor.Close()
		subscriberRedisConstructor.Close()
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
