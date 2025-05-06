package api

import (
	"backend/chat/usecase"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
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
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Image   string    `json:"image"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}

//type MessageHandler struct {
//	Pub pubsub.EventPublisher
//}

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

	if err := con.chatUseCase.CreateChat(newChat.Name, "myroom", newChat.Message); err != nil {
		log.Println("Error creating chat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
		return
	}

	//newChat.ID = uuid.New().String()
	//newChat.Date = time.Now()
	//
	//// newChat を byte 配列に変換
	//chatData, err := json.Marshal(newChat)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal chat"})
	//	return
	//}
	//
	//if err := p.Pub.Publish("chat_messages", newChat.ID, chatData); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish"})
	//	return
	//}

	c.JSON(http.StatusOK, newChat)
}
