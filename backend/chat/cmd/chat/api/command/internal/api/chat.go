package api

import (
	"backend/chat/cmd/chat/api/command/internal/usecase"
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
		chatGroup.POST("/:roomID/", con.postChat)
		chatGroup.POST("/:roomID/:chatId", con.editChat)
		chatGroup.DELETE("/:roomID/:chatId", con.deleteChat)
	}
}

func (con *ChatCommandController) postChat(ctx *gin.Context) {
	var newChat ChatRequest
	if err := ctx.ShouldBindJSON(&newChat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	senderID, err := uuid.Parse(newChat.SenderID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "sender must be a valid UUID"})
		return
	}

	roomID, err := uuid.Parse(ctx.Param("roomID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "roomID must be a valid UUID"})
	}

	if err := con.chatUseCase.CreateChat(senderID, roomID, newChat.Message); err != nil {
		log.Println("Error creating chat:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
		return
	}
	ctx.JSON(http.StatusOK, newChat)
}

func (con *ChatCommandController) editChat(ctx *gin.Context) {
	var editChat ChatRequest
	if err := ctx.ShouldBindJSON(&editChat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomID, err := uuid.Parse(ctx.Param("roomID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "roomID must be a valid UUID"})
		return
	}

	chatId, err := uuid.Parse(ctx.Param("chatId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "chatId must be a valid UUID"})
		return
	}

	if err := con.chatUseCase.EditChat(roomID, chatId, editChat.Message); err != nil {
		log.Println("Error updating chat:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to edit chat"})
		return
	}
	ctx.JSON(http.StatusOK, editChat)
}

func (con *ChatCommandController) deleteChat(ctx *gin.Context) {
	roomID, err := uuid.Parse(ctx.Param("roomID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "roomID must be a valid UUID"})
		return
	}

	chatId, err := uuid.Parse(ctx.Param("chatId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "chatId must be a valid UUID"})
		return
	}

	if err := con.chatUseCase.DeleteChat(roomID, chatId); err != nil {
		log.Println("Error deleting chat:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete chat"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "chat deleted successfully"})
}
