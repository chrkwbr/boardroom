package api

import (
	"backend/internal/chat/command/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	ID       string `json:"id"`
	SenderID string `json:"sender"`
	Message  string `json:"message"`
}

func (con *ChatCommandController) RegisterRoutes(r *gin.RouterGroup) {
	chatGroup := r.Group("/chats")
	{
		chatGroup.POST("/:roomId/", con.postChat)
		chatGroup.POST("/:roomId/:chatId", con.editChat)
		chatGroup.DELETE("/:roomId/:chatId", con.deleteChat)
	}
}

func (con *ChatCommandController) postChat(c *gin.Context) {
	var newChat ChatRequest
	if err := c.ShouldBindJSON(&newChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	senderID, err := uuid.Parse(newChat.SenderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sender must be a valid UUID"})
		return
	}

	roomID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roomId must be a valid UUID"})
	}

	if err := con.chatUseCase.CreateChat(senderID, roomID, newChat.Message); err != nil {
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
	chatId, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := con.chatUseCase.EditChat(chatId, editChat.Message); err != nil {
		log.Println("Error creating chat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
		return
	}
	c.JSON(http.StatusOK, editChat)
}

func (con *ChatCommandController) deleteChat(context *gin.Context) {
	chatId, err := uuid.Parse(context.Param("chatId"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := con.chatUseCase.DeleteChat(chatId); err != nil {
		log.Println("Error deleting chat:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete chat"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "chat deleted successfully"})
}
