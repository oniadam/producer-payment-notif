package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"producer-payment-notif/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func PublisherPaymentNotificationWaArray(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		res := models.Respons{
			Errors:            "1",
			ResponseCode:      "400",
			ResponseMessage:   "Error, Unmarshall body Request",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Data:              nil,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	req := []models.NotifPaymentWa{}

	if len(body) > 0 && body[0] == '[' {
		// jika array
		err = json.Unmarshal(body, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid array format"})
			return
		}
	} else {
		// jika single object
		var single models.NotifPaymentWa
		err = json.Unmarshal(body, &single)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid object format"})
			return
		}
		req = append(req, single)
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

	// data := models.NotifPaymentWa{}
	// for _, v := range req {
	// 	data = models.NotifPaymentWa{
	// 		AggrNo:            v.AggrNo,
	// 		TotalPaid:         v.TotalPaid,
	// 		IdNo:              v.IdNo,
	// 		WaNo:              v.WaNo,
	// 		CustomerName:      v.CustomerName,
	// 		Senddtm:           v.Senddtm,
	// 		Sendby:            v.Sendby,
	// 		Templatecode:      v.Templatecode,
	// 		TransactionSrc:    v.TransactionSrc,
	// 		Paymentmetodcode:  v.Paymentmetodcode,
	// 		Refno:             v.Refno,
	// 		RefNoWa:           v.RefNoWa,
	// 		Filepath:          v.Filepath,
	// 		Flagreversal:      v.Flagreversal,
	// 		Createdby:         v.Createdby,
	// 		Createddtm:        v.Createddtm,
	// 		CustomerServiceNo: v.CustomerServiceNo,
	// 		LanguageCode:      v.LanguageCode,
	// 	}
	// }

	// req = append(req, data)

	reqLog, _ := json.Marshal(req)

	traceID := uuid.New().String()

	// Publish pesan ke antrian
	errPub := ch.Publish(
		"payment_exchange", // exchange
		"payment.wa",       // routing key
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         reqLog,
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

	log.Println("Pesan berhasil dikirim:", string(reqLog))

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
