package api

import (
	"backend/chat/query"
	"backend/chat/query/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	message, err := con.chatService.ListMessage(ctx, room)
	if err != nil {
		log.Println("Error listing messages:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		return
	}
	var chatResponses []query.ChatResponse
	for _, chat := range message {
		chatResponse := query.ChatResponse{
			ID:        chat.ID,
			Sender:    chat.Sender,
			Message:   chat.Message,
			Version:   chat.Version,
			Timestamp: chat.CreatedAt,
		}
		chatResponses = append(chatResponses, chatResponse)
	}

	ctx.JSON(http.StatusOK, chatResponses)
}

func (con *ChatQueryController) history(ctx *gin.Context) {
	room := ctx.Param("room")
	chatId := ctx.Param("chatId")
	message, err := con.chatService.GetHistory(ctx, room, chatId)
	if err != nil {
		log.Println("Error listing messages:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		return
	}
	var chatResponses []query.ChatResponse
	for _, chat := range message {
		chatResponse := query.ChatResponse{
			ID:        chat.ID,
			Sender:    chat.Sender,
			Message:   chat.Message,
			Version:   chat.Version,
			Timestamp: chat.UpdatedAt,
		}
		chatResponses = append(chatResponses, chatResponse)
	}
	ctx.JSON(http.StatusOK, chatResponses)
}
