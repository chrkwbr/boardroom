package chat

import (
	"backend/event"
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

var initData = []Chat{}

func RegisterRoutes(r *gin.RouterGroup, hub *event.Hub) {
	chatGroup := r.Group("/chats")
	{
		chatGroup.GET("/:channelId/", getChats)
		chatGroup.POST("/:channelId/", func(c *gin.Context) {
			postChat(c, hub)
		})
	}
}

func getChats(c *gin.Context) {
	c.JSON(http.StatusOK, initData)
}

func postChat(c *gin.Context, hub *event.Hub) {
	var newChat Chat
	if err := c.ShouldBindJSON(&newChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newChat.ID = uuid.New().String()
	newChat.Date = time.Now()
	initData = append(initData, newChat)

	hub.BroadcastMessage(newChat)

	c.JSON(http.StatusOK, newChat)
}
