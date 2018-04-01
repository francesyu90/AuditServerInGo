package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"../data"
	"../exception"
	"../service"
	"../util"
)

type (
	Controller struct {
		service *service.Service
	}
)

func GetController(u *util.Utilities) *Controller {
	return newController(u)
}

func (controller Controller) Testing(c *gin.Context) {

	var acctTxnEvent data.QuoteServerEvent

	ch := controller.service.ProcessReqBody(c, &acctTxnEvent)
	resp := <-ch

	if resp != nil {
		handleError(exception.AS00001, resp, c)
	} else {
		c.JSON(http.StatusOK, acctTxnEvent)
	}

}

/*
	Private methods
*/

func newController(u *util.Utilities) *Controller {
	service := service.GetService(u)
	return &Controller{service}
}

func handleError(uuid exception.UUID, sysErr error, c *gin.Context) {
	err1 := exception.NewASError(uuid, sysErr.Error(), "")
	c.JSON(http.StatusInternalServerError, err1)
	log.Panicln(err1)
}
