package controller

import (
	"io/ioutil"
	"log"
	"net/http"

	"../data"
	"../exception"
	"../util"
	"github.com/gin-gonic/gin"
)

type (
	Controller struct {
		u *util.Utilities
	}
)

func GetController(u *util.Utilities) *Controller {
	return newController(u)
}

func (controller Controller) Testing(c *gin.Context) {

	var acctTxnEvent data.QuoteServerEvent

	ch := processReqBody(c, &acctTxnEvent)
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
	return &Controller{u}
}

func processReqBody(c *gin.Context, i interface{}) <-chan error {

	ch := make(chan error)

	go func() {

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			ch <- err
		}

		log.Println(string(body))
		util.UnserializeObject(body, &i)
		ch <- nil

	}()

	return ch

}

func handleError(uuid exception.UUID, sysErr error, c *gin.Context) {
	err1 := exception.NewASError(uuid, sysErr.Error(), "")
	c.JSON(http.StatusInternalServerError, err1)
	log.Panicln(err1)
}
