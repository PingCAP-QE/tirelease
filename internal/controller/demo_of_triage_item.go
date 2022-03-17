package controller

import (
	"net/http"

	"tirelease/internal/entity"
	"tirelease/internal/repository"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TriageOption struct {
	Owner string `json:"owner" form:"owner"`
	Repo  string `json:"repo" form:"repo"`
}

type TriageOprate struct {
	Owner  string   `json:"owner" form:"owner"`
	Repo   string   `json:"repo" form:"repo"`
	Number int      `json:"number" form:"number"`
	Lables []string `json:"labels" form:"labels"`
}

// Rest-API controller
func InsertTriageItems(c *gin.Context) {
	// Params
	triageOption := TriageOption{}
	if err := c.ShouldBindWith(&triageOption, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	triageItems, err := service.CollectTriageItemByRepo(triageOption.Owner, triageOption.Repo)
	if nil != err {
		c.Error(err)
		return
	}

	err2 := service.SavaTriageItems(triageItems)
	if nil != err2 {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func SelectTriageItems(c *gin.Context) {
	// Params
	option := entity.TriageItemOption{}
	if err := c.ShouldBindWith(&option, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	triageItems, err := repository.TriageItemSelect(&option)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": triageItems})
}

func AddLabelsToIssue(c *gin.Context) {
	// Params
	operate := TriageOprate{}
	if err := c.ShouldBindWith(&operate, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	err := service.AddLabelOfAccept(operate.Owner, operate.Repo, operate.Number, operate.Lables)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
