package controller

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"

	"../data"
	"../exception"
	services "../service"
	"../util"
)

type (
	Controller struct {
		service *services.Service
		session *mgo.Session
	}
)

func GetController(u *util.Utilities) *Controller {
	return newController(u)
}

func (controller Controller) Testing(c *gin.Context) {

	var acctTxnEvent data.QuoteServerEvent

	ch := controller.service.ProcessReqBody(c, &acctTxnEvent)
	resp := <-ch

	log.Println(controller.session)

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

	service := services.GetService(u)
	configService := services.GetConfigService(u)
	session := configService.GetDBConn()

	return &Controller{service, session}
}

func handleError(uuid exception.UUID, sysErr error, c *gin.Context) {
	err1 := exception.NewASError(uuid, sysErr.Error(), "")
	c.JSON(http.StatusInternalServerError, err1)
	log.Panicln(err1)
}
