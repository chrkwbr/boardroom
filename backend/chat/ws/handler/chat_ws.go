package handler

import (
	"backend/chat/domain"
	"backend/infra/hub"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatWebSocket struct {
}

func NewChatWebSocket() *ChatWebSocket {
	return &ChatWebSocket{}
}

type WsChatEvent struct {
	ChatId    string `json:"id"`
	EventType string `json:"event_type"`
	Sender    string `json:"sender"`
	Room      string `json:"room"`
	Message   string `json:"message"`
	Version   int64  `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // ToDo check origin
	},
}

func (ws *ChatWebSocket) RegisterRoutes(r *gin.RouterGroup) {
	chatGroup := r.Group("/")
	{
		chatGroup.GET("", func(c *gin.Context) {
			ws.handleWebSocketChat(c)
		})
	}
}

var activeSockets = struct {
	connections map[*websocket.Conn]bool
	mu          sync.Mutex
}{
	connections: make(map[*websocket.Conn]bool),
}

func (ws *ChatWebSocket) handleWebSocketChat(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}

	// 接続を追跡
	activeSockets.mu.Lock()
	activeSockets.connections[conn] = true
	activeSockets.mu.Unlock()

	wsChatEventPusherHub, err := hub.GetHubFactory().GetHub(hub.ChatEventWsPusher)
	if err != nil {
		log.Println("Failed to get hub:", err)
		err := conn.Close()
		if err != nil {
			return
		}
	}
	client := wsChatEventPusherHub.CreateAndRegisterClient(32)

	closeClient := func() {
		activeSockets.mu.Lock()
		delete(activeSockets.connections, conn)
		activeSockets.mu.Unlock()

		wsChatEventPusherHub.UnregisterClient(client)
		err := conn.Close()
		if err != nil {
			return
		}
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
		var chatEvent domain.ChatEvent
		if err := json.Unmarshal(msg, &chatEvent); err != nil {
			log.Println("Failed to unmarshal chat:", err)
			return
		}

		var wsChatEvent = &WsChatEvent{}
		fmt.Printf("Received chat event: %+v\n", chatEvent)

		switch chatEvent.Type {
		case domain.EventTypeCreated:
			var payload domain.ChatCreatedPayload
			if err := json.Unmarshal(chatEvent.Payload, &payload); err != nil {
				log.Println("Failed to unmarshal chat:", err)
				return
			}
			wsChatEvent = &WsChatEvent{
				ChatId:    payload.ID.String(),
				EventType: string(chatEvent.Type),
				Sender:    payload.SenderID.String(),
				Room:      payload.RoomID.String(),
				Timestamp: chatEvent.OccurredAt,
			}
		default:
			log.Println("Unknown event type:", chatEvent.Type)
			return
		}

		if err := conn.WriteJSON(wsChatEvent); err != nil {
			log.Println("WebSocket write error:", err)
			return
		}
	})

	conn.SetPingHandler(func(string) error {
		return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
	})
}
