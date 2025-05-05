package chat

import (
	"backend/event"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Chat struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Image   string    `json:"image"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}

type MessageHandler struct {
	Pub event.EventPublisher
}

var initData = []Chat{}

func (p *MessageHandler) RegisterRoutes(r *gin.RouterGroup) {
	chatGroup := r.Group("/chats")
	{
		chatGroup.GET("/:channelId/", getChats)
		chatGroup.POST("/:channelId/", p.postChat)
	}
}

func getChats(c *gin.Context) {
	c.JSON(http.StatusOK, initData)
}

func (p *MessageHandler) postChat(c *gin.Context) {
	var newChat Chat
	if err := c.ShouldBindJSON(&newChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newChat.ID = uuid.New().String()
	newChat.Date = time.Now()

	// newChat を byte 配列に変換
	chatData, err := json.Marshal(newChat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal chat"})
		return
	}

	if err := p.Pub.Publish("chat_messages", newChat.ID, chatData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish"})
		return
	}

	c.JSON(http.StatusOK, newChat)
}
