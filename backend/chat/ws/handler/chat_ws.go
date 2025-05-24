package handler

import (
	"backend/chat/command/domain"
	"backend/chat/event"
	"backend/infra/hub"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type ChatWebSocket struct {
}

func NewChatWebSocket() *ChatWebSocket {
	return &ChatWebSocket{}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // ToDo check origin
	},
}

func (ws *ChatWebSocket) RegisterRoutes(r *gin.RouterGroup) {
	chatGroup := r.Group("/chats")
	{
		chatGroup.GET("/:channelId/", func(c *gin.Context) {
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
		chatEvent := &event.ChatEvent{}
		if err := json.Unmarshal(msg, chatEvent); err != nil {
			log.Println("Failed to unmarshal chat:", err)
			return
		}
		chat := &domain.Chat{}
		if err := json.Unmarshal(chatEvent.Payload, chat); err != nil {
			log.Println("Failed to unmarshal chat:", err)
			return
		}

		wsChat := &WsChatEvent{
			ChatId:    chatEvent.ChatId.String(),
			EventType: chatEvent.EventType,
			Sender:    chat.Sender,
			Room:      chat.Room,
			Message:   chat.Message,
			Version:   chat.Version,
			Timestamp: chat.Timestamp,
		}

		if err := conn.WriteJSON(wsChat); err != nil {
			log.Println("WebSocket write error:", err)
			return
		}
	})

	conn.SetPingHandler(func(string) error {
		return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
	})
}
