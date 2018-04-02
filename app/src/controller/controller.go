package controller

import (
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
		u       *util.Utilities
	}
)

func GetController(u *util.Utilities) *Controller {
	return newController(u)
}

func (controller Controller) Testing(c *gin.Context) {

	ch := controller.checkAndHandleDBError(c)
	hasErr := <-ch
	if hasErr {
		return
	}

	events, asErr := controller.service.GetAllEvents()
	if asErr != nil {
		controller.handleError(c, asErr)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   events})

}

func (controller Controller) HandleQSEvent(c *gin.Context) {

	ch := controller.checkAndHandleDBError(c)
	hasErr := <-ch
	if hasErr {
		return
	}

	var qsEvent data.QuoteServerEvent
	ch1 := controller.service.ProcessReqBody(c, &qsEvent)
	asErr := <-ch1
	if asErr != nil {
		controller.handleError(c, asErr)
		return
	}

	event := &data.Event{
		EventType: data.QuSEvent,
		QsEvent:   &qsEvent,
	}

	asErr1 := controller.service.SaveEvent(event)
	if asErr1 != nil {
		controller.handleError(c, asErr1)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   qsEvent})

}

/*
	Private methods
*/

func newController(u *util.Utilities) *Controller {

	configService := services.GetConfigService(u)
	session, _ := configService.GetDBConn()
	service := services.GetService(u, session)

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
