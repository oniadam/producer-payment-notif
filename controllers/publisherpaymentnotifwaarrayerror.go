package controllers

import (
	"encoding/json"
	"log"
	"producer-payment-notif/repo"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func PublisherPaymentNotificationWaArrayError(c *gin.Context) {

	payload, _, _ := repo.GetDataQueueError()

	var data [][]map[string]interface{}

	err := json.Unmarshal([]byte(payload), &data)
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("Gagal terhubung ke RabbitMQ:", err)
		return
	}
	defer conn.Close()

	// Membuka channel
	ch, err := conn.Channel()
	if err != nil {
		log.Println("Gagal membuka channel", err)
		return
	}
	defer ch.Close()

	for _, arr := range data {
		body, err := json.Marshal(arr) // <-- tetap array
		if err != nil {
			log.Println("marshal error:", err)
			continue
		}

		err = ch.ExchangeDeclare(
			"payment_exchange",
			"topic",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Println("Gagal mendeklarasi antrian:", err)
			return
		}

		err = ch.ExchangeDeclare(
			"payment_dlx",
			"topic",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatal("Gagal declare DLX:", err)
		}

		args := amqp.Table{
			"x-dead-letter-exchange":    "payment_dlx",
			"x-dead-letter-routing-key": "payment.error",
		}

		_, err = ch.QueueDeclare(
			"paymentnotificationwa_queue",
			true,
			false,
			false,
			false,
			args,
		)
		if err != nil {
			log.Fatal("Gagal declare queue:", err)
		}

		// Bind queue ke exchange
		err = ch.QueueBind(
			"paymentnotificationwa_queue",
			"payment.wa",
			"payment_exchange",
			false,
			nil,
		)
		if err != nil {
			log.Fatal("Gagal bind queue:", err)
		}

		// ================================
		// 4. Declare Error Queue
		// ================================
		_, err = ch.QueueDeclare(
			"paymentnotificationwa_error_queue",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatal("Gagal declare error queue:", err)
		}

		err = ch.QueueBind(
			"paymentnotificationwa_error_queue",
			"payment.error",
			"payment_dlx",
			false,
			nil,
		)
		if err != nil {
			log.Fatal("Gagal bind error queue:", err)
		}

		traceID := uuid.New().String()

		// Publish pesan ke antrian
		errPub := ch.Publish(
			"payment_exchange", // exchange
			"payment.wa",       // routing key
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: amqp.Persistent,
				Headers: amqp.Table{
					"trace_id": traceID,
					"service":  "payment_notif_wa",
				},
			},
		)
		if errPub != nil {
			log.Println("Gagal mengirim pesan:", err)
			return
		}

		log.Println("Pesan berhasil dikirim:", string(body))

	}

	// c.String(http.StatusOK, "Pesan berhasil dikirim ke RabbitMQ")
	c.JSON(200, gin.H{
		"error":             "",
		"responseCode":      "200",
		"responseMessage":   "success",
		"responseTimestamp": time.Now().Format("2006-01-02 15:04:05"),
		"data":              nil,
	})
	// return
}
