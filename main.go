package main

import (
	"log"
	"producer-payment-notif/controllers"
	"producer-payment-notif/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	g := gin.Default()

	g.Use(utils.RequestLogger())

	g.POST("/api/v1/send-msg-notif", controllers.PublisherPaymentNotification)
	g.POST("/api/v1/send-msg-notif-wa", controllers.PublisherPaymentNotificationWa)
	g.GET("/api/v1/list/queue", controllers.ListQueue)
	g.GET("/api/v1/list/queue/:name", controllers.DetailQueue)
	g.DELETE("/api/v1/delete/queue/:name", controllers.DeleteQueue)

	g.Run(":8913")
}
