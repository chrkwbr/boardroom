package api

import (
	"backend/chat/command/usecase"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ChatCommandController struct {
	chatUseCase *usecase.ChatUseCase
}

func NewChatCommandController(chatUseCase *usecase.ChatUseCase) *ChatCommandController {
	return &ChatCommandController{
		chatUseCase: chatUseCase,
	}
}

type ChatRequest struct {
	ID      string `json:"id"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

func (con *ChatCommandController) RegisterRoutes(r *gin.RouterGroup) {
	chatGroup := r.Group("/chats")
	{
		chatGroup.POST("/:channelId/", con.postChat)
	}
}

func (con *ChatCommandController) postChat(c *gin.Context) {
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
