package main

import (
	"log"
	"producer-payment-notif/controllers"
	"producer-payment-notif/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	passphrase := "mNsRjOIdbyj1X2i6lLFJ5KE/evhYQIEz"

	// Muat dan dekripsi file .env.enc
	err := utils.LoadEncryptedEnv(".env.enc", passphrase)
	if err != nil {
		log.Fatal("Error loading encrypted .env file:", err)
	}

	g := gin.Default()

	g.Use(utils.RequestLogger())

	g.POST("/api/v1/send-msg-notif", controllers.PublisherPaymentNotification)
	g.POST("/api/v1/send-msg-notif-wa", controllers.PublisherPaymentNotificationWaArray)
	g.GET("/api/v1/list/queue", controllers.ListQueue)
	g.GET("/api/v1/monitor/queue/error", controllers.MonitorQueueError)
	g.GET("/api/v1/list/queue/:name", controllers.DetailQueue)
	g.DELETE("/api/v1/delete/queue/:name", controllers.DeleteQueue)

	g.POST("/api/v1/enkripenv", controllers.EnkripEnv)
	g.POST("/api/v1/dekripenv", controllers.DekripEnv)

	g.Run(":8913")
}
