package controller

import (
	"log"
	"net/http"

	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Rest-API controller
func WebhookHandler(c *gin.Context) {
	webhookPayload := service.WebhookPayload{}
	if err := c.ShouldBindWith(&webhookPayload, binding.JSON); err != nil {
		log.Fatal(err)
		c.Error(err)
		return
	}
	if err := service.UpdatePrAndIssue(webhookPayload); err != nil {
		log.Fatal(err)
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func InitDataForDemo(c *gin.Context) {
	service.InitDB()
}
