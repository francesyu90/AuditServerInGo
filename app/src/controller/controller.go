package controller

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"

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

	log.Println("Hello World")

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})

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
