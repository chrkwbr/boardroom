package chat

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Chat struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Image   string    `json:"image"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}

var initData = []Chat{}

func RegisterRoutes(r *gin.Engine) {
	chatGroup := r.Group("/chats")
	{
		chatGroup.GET("/:channelId/", getChats)
		chatGroup.POST("/:channelId/", postChat)
	}
}

func getChats(c *gin.Context) {
	c.JSON(http.StatusOK, initData)
}

func postChat(c *gin.Context) {
	var newChat Chat
	if err := c.ShouldBindJSON(&newChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newChat.ID = len(initData)
	newChat.Date = time.Now()
	initData = append(initData, newChat)
	c.JSON(http.StatusOK, newChat)
}
