package api

import (
	"backend/chat/query/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
		chatGroup.GET("/:room/", con.list)
		chatGroup.GET("/:room/:chatId/history/", con.history)
	}
}

func (con *ChatQueryController) list(ctx *gin.Context) {
	room := ctx.Param("room")
	chats, err := con.chatService.ListMessage(ctx, room)
	if err != nil {
		log.Println("Error listing messages:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		return
	}
	ctx.JSON(http.StatusOK, chats)
}

func (con *ChatQueryController) history(ctx *gin.Context) {
	room := ctx.Param("room")
	chatId := ctx.Param("chatId")
	chats, err := con.chatService.GetHistory(ctx, room, chatId)
	if err != nil {
		log.Println("Error listing messages:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		return
	}
	ctx.JSON(http.StatusOK, chats)
}
