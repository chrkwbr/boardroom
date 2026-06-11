package main

import (
	chatCommandApi "backend/internal/chat/command/api"
	"backend/internal/chat/command/usecase"
	chatqueryApi "backend/internal/chat/query/api"
	"backend/internal/chat/query/repository"
	"backend/internal/chat/query/service"
	"backend/internal/chat/readmodel"
	wschat "backend/internal/chat/ws/handler"
	"backend/internal/chat/ws/processor"
	"backend/internal/shared/infra/pubsub/kafka"
	"database/sql"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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
	//chatOutboxProcessor := processor2.NewOutboxProcessor(
	//	ChatDB,
	//	persistence2.NewChatOutboxRepositoryImpl(),
	//	eventPublisher)

	chatCommandApi := chatCommandApi.NewChatCommandController(
		usecase.NewChatUseCase(
			eventPublisher,
			//persistence2.NewChatRepositoryImpl(),
			//persistence2.NewChatOutboxRepositoryImpl(),
			//tx.NewTransactionManager(ChatDB),
		),
	)
	scylla, err := readmodel.NewChatScyllaRepository("localhost")
	if err != nil {
		log.Fatal("Failed to connect to ScyllaDB:", err)
	}
	chatQueryApi := chatqueryApi.NewChatQueryController(
		service.NewChatService(repository.NewChatScyllaQueryRepository(scylla)),
	)

	api := r.Group("/api")
	chatCommandApi.RegisterRoutes(api)
	chatQueryApi.RegisterRoutes(api)

	//apinotification.RegisterRoutes(api)

	// == Chat WebSocket ==
	//subscriberWebsocketProcessor := kafka.NewKafkaReader([]string{"localhost:9092"}, "chat-events", "websocket_processor")
	processor.NewChatRedisSubscriber(RedisClient).Start()

	ws := r.Group("/ws")
	chatWs := wschat.NewChatWebSocket()
	chatWs.RegisterRoutes(ws)

	subscriberRedisConstructor := kafka.NewKafkaReader([]string{"localhost:9092"}, "chat-events", "redis_constructor")
	readmodelRedis := readmodel.NewChatRedisRepository(RedisClient)
	readmodel.NewRedisConstructor(subscriberRedisConstructor, readmodelRedis).Start()

	defer func() {
		if err := ChatDB.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
		//chatOutboxProcessor.Close()

		if err := RedisClient.Close(); err != nil {
			log.Println("Error closing Redis client:", err)
		}

		eventPublisher.Close()
		//subscriberWebsocketProcessor.Close()
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

	cdb, err := sql.Open("postgres", "host=localhost port=5433 user=boardroom password=boardroom dbname=boardroom search_path=chat sslmode=disable")
	if err != nil {
		panic(err)
	}
	ChatDB = cdb
}
