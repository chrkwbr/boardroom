package api

import (
	"backend/chat/command/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		chatGroup.POST("/:roomId/", con.postChat)
		chatGroup.POST("/:roomId/:chatId", con.editChat)
	}
}

func (con *ChatCommandController) postChat(c *gin.Context) {
	var newChat ChatRequest
	if err := c.ShouldBindJSON(&newChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := con.chatUseCase.CreateChat(newChat.Sender, c.Param("roomId"), newChat.Message); err != nil {
		log.Println("Error creating chat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
		return
	}
	c.JSON(http.StatusOK, newChat)
}

func (con *ChatCommandController) editChat(c *gin.Context) {
	var editChat ChatRequest
	if err := c.ShouldBindJSON(&editChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cahtId, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := con.chatUseCase.EditChat(cahtId, editChat.Sender, c.Param("roomId"), editChat.Message); err != nil {
		log.Println("Error creating chat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
		return
	}
	c.JSON(http.StatusOK, editChat)
}
