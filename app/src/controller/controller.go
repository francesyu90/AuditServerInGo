package controller

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"

	"../exception"
	services "../service"
	"../util"
)

type (
	Controller struct {
		service *services.Service
		session *mgo.Session
		u       *util.Utilities
	}
)

func GetController(u *util.Utilities) *Controller {
	return newController(u)
}

func (controller Controller) Testing(c *gin.Context) {

	controller.session = nil

	ch := controller.checkAndHandleDBError(c)
	hasErr := <-ch
	if hasErr {
		return
	}

	log.Println("Hello World")

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})

}

func (controller Controller) HandleQSEvent(c *gin.Context) {

}

/*
	Private methods
*/

func newController(u *util.Utilities) *Controller {

	service := services.GetService(u)
	configService := services.GetConfigService(u)
	session, _ := configService.GetDBConn()

	return &Controller{service, session, u}
}

func (controller Controller) checkAndHandleDBError(
	c *gin.Context) <-chan bool {

	ch := make(chan bool)

	go func() {

		asErr := controller.checkDBConn()
		if asErr != nil {
			controller.handleError(c, asErr)
			ch <- true
		} else {
			ch <- false
		}

	}()

	return ch
}

func (controller Controller) handleError(
	c *gin.Context, asErr *exception.ASError) {

	c.JSON(http.StatusInternalServerError,
		gin.H{
			"status": http.StatusInternalServerError,
			"error":  asErr.ErrorMessage(),
		})
}

func (controller Controller) checkDBConn() *exception.ASError {

	if controller.session == nil {
		asError := controller.u.GetError(exception.AS00007, "db_conn_err", nil)
		return asError
	}
	return nil
}
