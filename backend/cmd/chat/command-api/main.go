package main

import (
	chatCommandApi "backend/internal/chat/command/api"
	"backend/internal/chat/command/usecase"
	chatqueryApi "backend/internal/chat/query/api"
	"backend/internal/chat/query/repository"
	"backend/internal/chat/query/service"
	"backend/internal/chat/ws/processor"
	"backend/internal/shared/infra/pubsub/kafka"
	"database/sql"
	"log"
	"os"

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
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// == Chat API ==
	eventPublisher := kafka.NewKafkaWriter([]string{"localhost:9092"})
	chatCommandApi := chatCommandApi.NewChatCommandController(
		usecase.NewChatUseCase(
			eventPublisher,
		),
	)
	chatReadModelRepository := repository.NewChatReadModelRepository(RedisClient)
	chatQueryApi := chatqueryApi.NewChatQueryController(
		service.NewChatService(chatReadModelRepository),
	)

	api := r.Group("/api")
	chatCommandApi.RegisterRoutes(api)
	chatQueryApi.RegisterRoutes(api)

	processor.NewChatRedisSubscriber(RedisClient).Start()

	defer func() {
		if err := ChatDB.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
		if err := RedisClient.Close(); err != nil {
			log.Println("Error closing Redis client:", err)
		}
		eventPublisher.Close()
	}()

	port := os.Getenv("PORT")
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	cdb, err := sql.Open("postgres", "host=localhost port=5433 user=boardroom password=boardroom dbname=boardroom search_path=chat sslmode=disable")
	if err != nil {
		panic(err)
	}
	ChatDB = cdb
}
