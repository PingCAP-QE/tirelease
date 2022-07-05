package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindUser(c *gin.Context) {
	// Params
	option := UserRequest{}
	if err := c.ShouldBindQuery(&option); err != nil {
		c.Error(err)
		return
	}

	// Action
	user, err := HandleFindUser(option)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
