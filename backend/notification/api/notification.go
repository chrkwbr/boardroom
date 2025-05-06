package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Notification struct {
	ID      int       `json:"id"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
	Read    bool      `json:"read"`
}

var notifications = []Notification{}

func RegisterRoutes(r *gin.RouterGroup) {
	notificationGroup := r.Group("/notifications")
	{
		notificationGroup.GET("/", getNotifications)
	}
}

func getNotifications(c *gin.Context) {
	n := make([]Notification, len(notifications))
	copy(n, notifications)
	notifications = notifications[:0]
	c.JSON(http.StatusOK, n)
}
