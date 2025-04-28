package chat

import (
	"backend/event"
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
	chatGroup := r.Group("/chats")
	{
		chatGroup.GET("/:channelId/", handleWebSocketChat)
	}
}

var activeSockets = struct {
	connections map[*websocket.Conn]bool
	mu          sync.Mutex
}{
	connections: make(map[*websocket.Conn]bool),
}

func handleWebSocketChat(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	// 接続を追跡
	activeSockets.mu.Lock()
	activeSockets.connections[conn] = true
	//log.Printf("WebSocket接続が確立されました。アクティブな接続数: %d", len(activeSockets.connections))
	activeSockets.mu.Unlock()

	defer func() {
		activeSockets.mu.Lock()
		delete(activeSockets.connections, conn)
		//log.Printf("WebSocket接続が終了しました。アクティブな接続数: %d", len(activeSockets.connections))
		activeSockets.mu.Unlock()
	}()

	conn.SetPingHandler(func(string) error {
		return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
	})

	eventCh := event.Subscribe("chat_created")
	defer func() {
		event.Unsubscribe("chat_created", eventCh)
	}()

	done := make(chan struct{})

	go func() {
		defer close(done)
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

	for {
		select {
		case evt, ok := <-eventCh:
			if !ok {
				return
			}
			err := conn.WriteJSON(evt)
			if err != nil {
				log.Println("WebSocket write error:", err)
				return
			}
		case <-done:
			// クライアントが切断した
			log.Println("WebSocket connection closed")
			return
		}
	}
}
