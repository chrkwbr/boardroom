package ws

import (
	"backend/chat/domain"
	"backend/infra/kafka"
	"backend/infra/pubsub"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // ToDo check origin
	},
}

func RegisterRoutes(r *gin.RouterGroup) {
	hub := pubsub.NewHub()
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

	chatGroup := r.Group("/chats")
	{
		chatGroup.GET("/:channelId/", func(c *gin.Context) {
			handleWebSocketChat(c, hub)
		})
	}
}

var activeSockets = struct {
	connections map[*websocket.Conn]bool
	mu          sync.Mutex
}{
	connections: make(map[*websocket.Conn]bool),
}

func handleWebSocketChat(c *gin.Context, hub *pubsub.Hub) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}

	// 接続を追跡
	activeSockets.mu.Lock()
	activeSockets.connections[conn] = true
	activeSockets.mu.Unlock()

	client := pubsub.NewClient(256)
	hub.RegisterClient(client)

	closeClient := func() {
		activeSockets.mu.Lock()
		delete(activeSockets.connections, conn)
		activeSockets.mu.Unlock()

		hub.UnregisterClient(client)
		conn.Close()
	}

	go func() {
		defer closeClient()

		for {
			// WebSocketからメッセージを読み取り続ける
			// クライアントが切断すると、このループは終了する
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error: %v", err)
				}
				return
			}
		}
	}()

	go client.Receive(func(msg []byte) {
		chat := domain.Chat{}
		if err := json.Unmarshal(msg, chat); err != nil {
			log.Println("Failed to unmarshal chat:", err)
			return
		}

		if err := conn.WriteJSON(chat); err != nil {
			log.Println("WebSocket write error:", err)
			return
		}
	})

	conn.SetPingHandler(func(string) error {
		return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
	})
}
