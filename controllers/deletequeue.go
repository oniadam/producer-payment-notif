package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func DeleteQueue(c *gin.Context) {
	name := c.Param("name")
	url := "http://localhost:15672/api/queues/%2F/" + name
	client := resty.New()
	resp, err := client.R().
		SetBasicAuth("guest", "guest").
		Delete(url)

	if err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode(), "application/json", resp.Body())
}
