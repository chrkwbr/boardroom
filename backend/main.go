package main

import (
	apichat "backend/api/chat"
	apinotification "backend/api/notification"
	"backend/event"
	"backend/infra"
	wschat "backend/ws/chat"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	hub := event.NewHub()
	go hub.Run()

	sub := infra.NewKafkaReader([]string{"localhost:9092"}, "chat_messages")
	go func() {
		if err := sub.Subscribe("_", func(key string, value []byte) error {
			hub.BroadcastMessage(value)
			return nil
		}); err != nil {
			panic(err)
		}
	}()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 許可するオリジン
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")

	pub := infra.NewKafkaWriter([]string{"localhost:9092"})
	h := &apichat.MessageHandler{Pub: pub}
	h.RegisterRoutes(api, hub)

	apinotification.RegisterRoutes(api)

	ws := r.Group("/ws")
	wschat.RegisterRoutes(ws, hub)

	r.Run()
}
