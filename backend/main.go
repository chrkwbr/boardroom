package main

import (
	apichat "backend/api/chat"
	apinotification "backend/api/notification"
	"backend/infra"
	wschat "backend/ws/chat"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 許可するオリジン
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")

	pub := infra.NewKafkaWriter([]string{"localhost:9092"})
	h := &apichat.MessageHandler{Pub: pub}
	h.RegisterRoutes(api)

	apinotification.RegisterRoutes(api)

	ws := r.Group("/ws")
	wschat.RegisterRoutes(ws)

	r.Run()
}
