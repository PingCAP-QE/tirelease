package controller

import (
	"fmt"
	"net/http"

	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/pkg/errors"
)

// Rest-API controller
func TestEntityInsert(c *gin.Context) {
	testEntity := entity.TestEntity{}
	if err := c.ShouldBindWith(&testEntity, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	if err := repository.TestEntityInsert(&testEntity); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func TestEntitySelect(c *gin.Context) {
	option := entity.TestEntityOption{}
	if err := c.ShouldBindWith(&option, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	testEntities, err := repository.TestEntitySelect(&option)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": testEntities})
}

type TestPingPongStruct struct {
	Name  string `json:"name" form:"name" uri:"name"`
	Other string `json:"other" form:"other" uri:"other"`
}

func TestPingError(c *gin.Context) {
	c.Error(errors.New("ping error"))
}

func TestPingPongGet(c *gin.Context) {
	param := TestPingPongStruct{}
	if err := c.ShouldBindWith(&param, binding.Form); err != nil {
		c.Error(err)
		return
	}
	fmt.Println(param)

	c.JSON(http.StatusOK, gin.H{"data": param})
}

func TestPingPongPost(c *gin.Context) {
	// param := TestPingPongStruct{}
	var param TestPingPongStruct
	if err := c.ShouldBindUri(&param); err != nil {
		c.Error(err)
		return
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.Error(err)
		return
	}
	fmt.Println(param)

	c.JSON(http.StatusOK, gin.H{"data": param})
}
