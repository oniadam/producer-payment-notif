package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"producer-payment-notif/repo"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func MonitorQueueError(c *gin.Context) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	queueName := "paymentnotificationwa_error_queue"

	var allMessages []interface{}

	for {
		msg, ok, err := ch.Get(queueName, false)
		if err != nil {
			log.Fatal(err)
		}

		if !ok {
			break // tidak ada message lagi
		}

		var jsonBody interface{}
		err = json.Unmarshal(msg.Body, &jsonBody)
		if err != nil {
			// kalau bukan JSON, tampilkan raw string
			jsonBody = string(msg.Body)
		}

		allMessages = append(allMessages, jsonBody)

		// balikin lagi ke queue
		msg.Nack(false, false)
		// msg.Nack(false, false) kalo mau pesan ga ilang
	}

	_, errinsPaymentNotifWa := repo.InsertQueueError(allMessages)
	if errinsPaymentNotifWa != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errinsPaymentNotifWa.Error(),
		})
		return
	}

	log.Println("Total Error Message:", len(allMessages))
	for _, m := range allMessages {
		log.Println("Message:", m)
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(allMessages),
		"data":  allMessages,
	})
}
