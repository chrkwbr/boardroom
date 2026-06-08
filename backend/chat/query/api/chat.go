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

type ChatResponse struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	SenderIcon string `json:"sender_icon"`
	RoomID     string `json:"room_id"`
	Version    int64  `json:"version"`
	Timestamp  int64  `json:"date"`
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
	var chatResponses []ChatResponse
	for _, chat := range message {
		chatResponse := ChatResponse{
			ID:         chat.ID,
			Message:    chat.Message,
			SenderID:   chat.Sender.ID.String(),
			SenderName: chat.Sender.Name,
			SenderIcon: chat.Sender.Icon,
			RoomID:     chat.RoomID,
			Version:    chat.Version,
			Timestamp:  chat.CreatedAt,
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
	var chatResponses []ChatResponse
	for _, chat := range message {
		chatResponse := ChatResponse{
			ID:         chat.ID,
			Message:    chat.Message,
			SenderID:   chat.Sender.ID.String(),
			SenderName: chat.Sender.Name,
			SenderIcon: chat.Sender.Icon,
			RoomID:     chat.RoomID,
			Version:    chat.Version,
			Timestamp:  chat.UpdatedAt,
		}
		chatResponses = append(chatResponses, chatResponse)
	}
	ctx.JSON(http.StatusOK, chatResponses)
}
