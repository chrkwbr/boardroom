package handler

import (
	"boardroom/chat-notification"
	hub "boardroom/shared/event-hub"
	"encoding/json"
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
	EventType notification.EventType `json:"event_type"`
	RoomId    string                 `json:"room_id"`
	ChatId    string                 `json:"chat_id"`
	Chat      notification.Chat      `json:"chat"`
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

	go client.Receive(receive(conn))

	conn.SetPingHandler(func(string) error {
		return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
	})
}

func receive(conn *websocket.Conn) func([]byte) {
	return func(msg []byte) {
		e := &notification.ChatRedisEvent{}
		if err := json.Unmarshal(msg, e); err != nil {
			log.Println("Unmarshal err:", err)
		}

		var chatEvent = &WsChatEvent{}

		switch e.Type {
		case notification.EventTypeCreated:
			chatEvent = &WsChatEvent{
				EventType: e.Type,
				RoomId:    e.Payload.RoomID.String(),
				ChatId:    e.Payload.ID.String(),
				Chat:      *e.Payload,
			}
		case notification.EventTypeUpdated:
			chatEvent = &WsChatEvent{
				EventType: e.Type,
				RoomId:    e.Payload.RoomID.String(),
				ChatId:    e.Payload.ID.String(),
				Chat:      *e.Payload,
			}
		case notification.EventTypeDeleted:
			chatEvent = &WsChatEvent{
				EventType: e.Type,
				RoomId:    e.RoomID.String(),
				ChatId:    e.ChatID.String(),
			}
		default:
			log.Println("Unknown event type:", e.Type)
			return
		}

		if err := conn.WriteJSON(chatEvent); err != nil {
			log.Println("WebSocket write error:", err)
			return
		}
	}
}
