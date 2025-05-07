package api

import (
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

type ChatResponse struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Image     string `json:"image"`
	Message   string `json:"message"`
	Timestamp int64  `json:"date"`
}

func (con *ChatQueryController) RegisterRoutes(r *gin.RouterGroup) {
	chatGroup := r.Group("/chats")
	{
		chatGroup.GET("/:room/", con.list)
	}
}

func (con *ChatQueryController) list(c *gin.Context) {
	room := c.Param("room")
	message, err := con.chatService.ListMessage(c, room)
	if err != nil {
		log.Println("Error listing messages:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		return
	}
	// to ChatResponse
	var chatResponses []ChatResponse
	for _, chat := range message {
		chatResponse := ChatResponse{
			ID:        chat.ID.String(),
			Sender:    chat.Sender,
			Image:     "https://img.daisyui.com/images/profile/demo/1@94.webp", // Placeholder for image URL
			Message:   chat.Message,
			Timestamp: chat.Timestamp,
		}
		chatResponses = append(chatResponses, chatResponse)
	}

	c.JSON(http.StatusOK, chatResponses)
}
