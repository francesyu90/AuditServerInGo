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
		loggers *util.Logger
	}
)

func GetController(u *util.Utilities, loggers *util.Logger) *Controller {
	return newController(u, loggers)
}

func (controller Controller) Testing(c *gin.Context) {

	ch := controller.checkAndHandleDBError(c)
	hasErr := <-ch
	if hasErr {
		return
	}

	// events, asErr := controller.service.GetAllEvents()
	// if asErr != nil {
	// 	controller.handleError(c, asErr)
	// 	return
	// } else if events == nil {
	// 	controller.handleNoEventsAvail("no_events_found_warning", c)
	// 	return
	// }

	// controller.loggers.INFO.Println("events: ", events)

	// c.JSON(
	// 	http.StatusOK,
	// 	gin.H{
	// 		"status": http.StatusOK,
	// 		"data":   events})

	controller.service.LogAll()

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
		})

}

func (controller Controller) LogAll(c *gin.Context) {

	ch := controller.checkAndHandleDBError(c)
	hasErr := <-ch
	if hasErr {
		return
	}

	asErr, asWarning := controller.service.LogAll()
	if asErr != nil {
		controller.handleError(c, asErr)
		return
	} else if asWarning != nil {
		controller.handleWarning(c, asWarning)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
		})

}

func (controller Controller) LogByUser(c *gin.Context) {

	ch := controller.checkAndHandleDBError(c)
	hasErr := <-ch
	if hasErr {
		return
	}

	userID := c.Param("userName")
	asErr, asWarning := controller.service.LogByUser(userID)
	if asErr != nil {
		controller.handleError(c, asErr)
		return
	} else if asWarning != nil {
		controller.handleWarning(c, asWarning)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
		})

}

func (controller Controller) HandleATEvent(c *gin.Context) {

	ch := controller.checkAndHandleDBError(c)
	hasErr := <-ch
	if hasErr {
		return
	}

	var atEvent data.AccountTransactionEvent
	ch1 := controller.service.ProcessReqBody(c, &atEvent)
	asErr := <-ch1
	if asErr != nil {
		controller.handleError(c, asErr)
		return
	}

	controller.loggers.INFO.Println("Account Transaction Event: ", atEvent)

	event := &data.Event{
		EventType:    data.AcTxnEvent,
		AcctTxnEvent: &atEvent,
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
			"data":   atEvent})

}

func (controller Controller) HandleSEvent(c *gin.Context) {

	ch := controller.checkAndHandleDBError(c)
	hasErr := <-ch
	if hasErr {
		return
	}

	var sEvent data.SystemEvent
	ch1 := controller.service.ProcessReqBody(c, &sEvent)
	asErr := <-ch1
	if asErr != nil {
		controller.handleError(c, asErr)
		return
	}

	controller.loggers.INFO.Println("System Event: ", sEvent)

	event := &data.Event{
		EventType: data.SyEvent,
		SysEvent:  &sEvent,
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
			"data":   sEvent})

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

	controller.loggers.INFO.Println("Quote Server Event: ", qsEvent)

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

func (controller Controller) HandleEEvent(c *gin.Context) {

	ch := controller.checkAndHandleDBError(c)
	hasErr := <-ch
	if hasErr {
		return
	}

	var eEvent data.ErrorEvent
	ch1 := controller.service.ProcessReqBody(c, &eEvent)
	asErr := <-ch1
	if asErr != nil {
		controller.handleError(c, asErr)
		return
	}

	controller.loggers.INFO.Println("Error Event: ", eEvent)

	event := &data.Event{
		EventType: data.ErEvent,
		ErrEvent:  &eEvent,
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
			"data":   eEvent})

}

/*
	Private methods
*/

func newController(u *util.Utilities, loggers *util.Logger) *Controller {

	configService := services.GetConfigService(u, loggers)
	session, _ := configService.GetDBConn()
	service := services.GetService(u, session, loggers)

	return &Controller{service, session, u, loggers}
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

	controller.loggers.ERROR.Println(asErr.ErrorMessage())

	c.JSON(http.StatusInternalServerError,
		gin.H{
			"status": http.StatusInternalServerError,
			"error":  asErr.ErrorMessage(),
		})
}

func (controller Controller) handleWarning(
	c *gin.Context, asWarning *exception.ASWarning) {

	controller.loggers.WARNING.Println(asWarning.WarningMessage())
	c.JSON(http.StatusNotFound,
		gin.H{
			"status":  http.StatusNotFound,
			"warning": asWarning.WarningMessage(),
		})
}

func (controller Controller) checkDBConn() *exception.ASError {

	if controller.session == nil {
		asError := controller.u.GetError(exception.AS00007, "db_conn_err", nil)
		return asError
	}
	return nil
}
