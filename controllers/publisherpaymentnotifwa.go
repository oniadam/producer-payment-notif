package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"producer-payment-notif/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func PublisherPaymentNotificationWa(c *gin.Context) {
	req := models.NotifPaymentWa{}
	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(dataBodyReq, &req)
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

	// Mendeklarasikan antrian (queue) yang akan digunakan
	queueName := "paymentnotificationwa_queue"
	_, err = ch.QueueDeclare(
		queueName, // Nama antrian
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Println("Gagal mendeklarasi antrian:", err)
		return
	}

	data := models.NotifPaymentWa{
		AggrNo:           req.AggrNo,
		TotalPaid:        req.TotalPaid,
		WaNo:             req.WaNo,
		CustomerName:     req.CustomerName,
		Senddtm:          req.Senddtm,
		Sendby:           req.Sendby,
		Templatecode:     req.Templatecode,
		TransactionSrc:   req.TransactionSrc,
		Paymentmetodcode: req.Paymentmetodcode,
		Refno:            req.Refno,
		RefNoWa:          req.RefNoWa,
		Filepath:         req.Filepath,
		Flagreversal:     req.Flagreversal,
		Createdby:        req.Createdby,
		Createddtm:       req.Createddtm,
	}
	reqLog, _ := json.Marshal(data)

	// Publish pesan ke antrian
	errPub := ch.Publish(
		"",        // exchange
		queueName, // routing key (nama antrian)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         reqLog, //[]byte(string(reqLog)),
			DeliveryMode: amqp.Persistent,
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
