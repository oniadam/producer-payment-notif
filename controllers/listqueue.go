package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func ListQueue(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetBasicAuth("guest", "guest").
		Get("http://localhost:15672/api/queues")

	if err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode(), "application/json", resp.Body())
}
