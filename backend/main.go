package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Chat struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Image   string    `json:"image"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}

var initData = []Chat{
	{
		ID:      0,
		Name:    "Dio Lupa",
		Image:   "https://img.daisyui.com/images/profile/demo/1@94.webp",
		Message: `"Remaining Reason" became an instant hit, praised for its haunting sound and emotional depth. A viral performance brought it widespread recognition, making it one of Dio Lupa’s most iconic tracks.`,
		Date:    time.Now(),
	},
	{
		ID:      1,
		Name:    "Ellie Beilish",
		Image:   "https://img.daisyui.com/images/profile/demo/4@94.webp",
		Message: `"Bears of a Fever" captivated audiences with its intense energy and mysterious lyrics. Its popularity skyrocketed after fans shared it widely online, earning Ellie critical acclaim.`,
		Date:    time.Now(),
	},
	{
		ID:      2,
		Name:    "Sabrino Gardener",
		Image:   "https://img.daisyui.com/images/profile/demo/3@94.webp",
		Message: `"Cappuccino" quickly gained attention for its smooth melody and relatable themes. The song’s success propelled Sabrino into the spotlight, solidifying their status as a rising star.`,
		Date:    time.Now(),
	},
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 許可するオリジン
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	r.GET("/:channelId/", func(c *gin.Context) {
		c.JSON(http.StatusOK, initData)
	})

	r.POST("/:channelId/", func(c *gin.Context) {
		var chat Chat
		if err := c.ShouldBindJSON(&chat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		chat.ID = len(initData)
		initData = append(initData, chat)
		c.JSON(http.StatusOK, chat)
	})

	r.Run()
}
