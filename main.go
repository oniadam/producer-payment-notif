package main

import (
	"producer-payment-notif/controllers"
	"producer-payment-notif/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

	g.Use(utils.RequestLogger())

	g.POST("/api/v1/send-msg-notif", controllers.PublisherPaymentNotification)
	g.POST("/api/v1/send-msg-notif-wa", controllers.PublisherPaymentNotificationWa)
	g.GET("/api/v1/list/queue", controllers.ListQueue)
	g.GET("/api/v1/list/queue/:name", controllers.DetailQueue)

	g.Run(":8913")
}
