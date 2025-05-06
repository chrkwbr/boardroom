package api

import (
	"backend/chat/usecase"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ChatController struct {
	chatUseCase *usecase.ChatUseCase
}

func NewChatController(chatUseCase *usecase.ChatUseCase) *ChatController {
	return &ChatController{
		chatUseCase: chatUseCase,
	}
}

type ChatRequest struct {
	ID      string `json:"id"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

func (con *ChatController) RegisterRoutes(r *gin.RouterGroup) {
	chatGroup := r.Group("/chats")
	{
		chatGroup.GET("/:channelId/", func(c *gin.Context) {
			c.JSON(http.StatusOK, []ChatRequest{})
		})
		chatGroup.POST("/:channelId/", con.postChat)
	}
}

func (con *ChatController) postChat(c *gin.Context) {
	var newChat ChatRequest
	if err := c.ShouldBindJSON(&newChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := con.chatUseCase.CreateChat(newChat.Sender, "myroom", newChat.Message); err != nil {
		log.Println("Error creating chat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
		return
	}
	c.JSON(http.StatusOK, newChat)
}
