package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func DetailQueue(c *gin.Context) {
	name := c.Param("name")

	client := resty.New()
	url := "http://localhost:15672/api/queues/%2F/" + name

	resp, err := client.R().
		SetBasicAuth("guest", "guest").
		Get(url)

	if err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp.StatusCode() == 404 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Queue not found"})
		return
	}

	c.Data(resp.StatusCode(), "application/json", resp.Body())
}
