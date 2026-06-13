package api

import (
	"backend/chat/cmd/api-query/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatQueryController struct {
	chatService *service.ChatService
}

func NewChatQueryController(chatService *service.ChatService) *ChatQueryController {
	return &ChatQueryController{
		chatService: chatService,
	}
}

func (con *ChatQueryController) RegisterRoutes(r *gin.RouterGroup) {
	chatGroup := r.Group("/chats")
	{
		chatGroup.GET("/:roomID/", con.list)
		chatGroup.GET("/:roomID/:chatID/history/", con.history)
	}
}

func (con *ChatQueryController) list(ctx *gin.Context) {
	r := ctx.Param("roomID")
	roomID, err := uuid.Parse(r)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "roomID must be a valid UUID"})
		return
	}

	chats, err := con.chatService.ListMessage(ctx, roomID)
	if err != nil {
		log.Println("Error listing messages:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		return
	}
	ctx.JSON(http.StatusOK, chats)
}

func (con *ChatQueryController) history(ctx *gin.Context) {
	r := ctx.Param("roomID")
	roomID, err := uuid.Parse(r)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "roomID must be a valid UUID"})
		return
	}
	c := ctx.Param("chatID")
	chatID, err := uuid.Parse(c)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "chatID must be a valid UUID"})
		return
	}

	chats, err := con.chatService.GetHistory(ctx, roomID, chatID)
	if err != nil {
		log.Println("Error listing messages:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		return
	}
	ctx.JSON(http.StatusOK, chats)
}
