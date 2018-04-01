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

	asErr := controller.checkDBConn()
	if asErr != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"status": http.StatusInternalServerError,
				"error":  asErr.ErrorMessage(),
			})
		return
	}

	log.Println("Hello World")

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})

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

func (controller Controller) checkDBConn() *exception.ASError {

	if controller.session == nil {
		asError := controller.u.GetError(exception.AS00007, "db_conn_err", nil)
		return asError
	}
	return nil
}
