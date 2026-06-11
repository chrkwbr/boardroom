package handler

import (
	"backend/internal/chat/readmodel"
	"backend/internal/shared/infra/hub"
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
	EventType readmodel.EventType     `json:"event_type"`
	RoomId    string                  `json:"room_id"`
	ChatId    string                  `json:"chat_id"`
	Chat      readmodel.ChatReadModel `json:"chat"`
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
		e := &readmodel.ChatRedisEvent{}
		if err := json.Unmarshal(msg, e); err != nil {
			log.Println("Unmarshal err:", err)
		}

		var wsChatEvent = &WsChatEvent{}

		switch e.Type {
		case readmodel.EventTypeCreated:
			wsChatEvent = &WsChatEvent{
				EventType: e.Type,
				RoomId:    e.Payload.RoomID.String(),
				ChatId:    e.Payload.ID.String(),
				Chat:      *e.Payload,
			}
		case readmodel.EventTypeUpdated:
			wsChatEvent = &WsChatEvent{
				EventType: e.Type,
				RoomId:    e.Payload.RoomID.String(),
				ChatId:    e.Payload.ID.String(),
				Chat:      *e.Payload,
			}
		case readmodel.EventTypeDeleted:
			wsChatEvent = &WsChatEvent{
				EventType: e.Type,
				RoomId:    e.Payload.RoomID.String(),
				ChatId:    e.Payload.ID.String(),
			}
		default:
			log.Println("Unknown event type:", e.Type)
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
