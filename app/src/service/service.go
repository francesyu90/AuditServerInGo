package service

import (
	"encoding/xml"
	"io/ioutil"

	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"

	"../data"
	"../exception"
	"../repository"
	"../util"
)

type Service struct {
	u       *util.Utilities
	session *mgo.Session
	repo    *repository.Repository
	loggers *util.Logger
}

func GetService(
	u *util.Utilities,
	session *mgo.Session,
	loggers *util.Logger) *Service {

	return newService(u, session, loggers)
}

func (service Service) ProcessReqBody(
	c *gin.Context, i interface{}) <-chan *exception.ASError {

	ch := make(chan *exception.ASError)

	go func() {

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			asErr := service.u.GetError(exception.AS00003, "", err)
			service.loggers.ERROR.Println(asErr.ErrorMessage())
			ch <- asErr
		}

		service.loggers.INFO.Println("Request Body: ", string(body))

		asErr1 := service.u.UnserializeObject(body, &i)
		if asErr1 != nil {
			ch <- asErr1
		} else {
			ch <- nil
		}

	}()

	return ch
}

func (service Service) SaveEvent(
	event *data.Event) *exception.ASError {

	return service.repo.Insert(event)
}

func (service Service) GetAllEvents() (
	[]*data.Event, *exception.ASError) {

	service.loggers.INFO.Println("Get all events")

	return service.repo.FindAll()
}

func (service Service) GetAllEventsByUser(
	userID string) ([]*data.Event, *exception.ASError) {

	service.loggers.INFO.Println("Get all events with user id: ", userID)
	return service.repo.FindByUserID(userID)
}

func (service Service) LogAll() (*exception.ASError, *exception.ASWarning) {

	events, asErr := service.GetAllEvents()
	if asErr != nil {
		return asErr, nil
	} else if events == nil {
		asWarning := service.u.GetWarning(exception.AS00014, "no_events_found_warning")
		return nil, asWarning
	}

	for _, event := range events {
		xmlString := service.getXMLEventString(event)
		service.loggers.XML.Println(xmlString)
	}

	return nil, nil
}

func (service Service) LogByUser(userID string) (
	*exception.ASError, *exception.ASWarning) {

	events, asErr := service.GetAllEventsByUser(userID)
	if asErr != nil {
		return asErr, nil
	} else if events == nil {
		asWarning :=
			service.u.GetWarning(
				exception.AS00016,
				"no_events_found_by_user_warning")
		return nil, asWarning
	}

	for _, event := range events {
		xmlString := service.getXMLEventString(event)
		service.loggers.XML.Println(xmlString)
	}

	return nil, nil
}

/*
	Private methods
*/

func newService(
	u *util.Utilities,
	session *mgo.Session,
	loggers *util.Logger) *Service {

	repo := repository.GetRepository(session, u, loggers)
	return &Service{u, session, repo, loggers}
}

func getEvent(event *data.Event) interface{} {

	switch event.EventType {
	case data.AcTxnEvent:
		atEvent := event.AcctTxnEvent
		return data.GetAccountTransaction(
			atEvent.Server,
			atEvent.TransactionNum,
			atEvent.Action,
			atEvent.UserId,
			atEvent.Funds,
			atEvent.Timestamp)
	case data.SyEvent:
		sysEvent := event.SysEvent
		return data.GetSystemEvent(
			sysEvent.Server,
			sysEvent.TransactionNum,
			sysEvent.Command,
			sysEvent.UserId,
			sysEvent.StockSymbol,
			sysEvent.Funds,
			sysEvent.Timestamp)
	case data.QuSEvent:
		qusEvent := event.QsEvent
		return data.GetQuoteServer(
			qusEvent.Server,
			qusEvent.TransactionNum,
			qusEvent.QuoteServerEventTime,
			qusEvent.Command,
			qusEvent.UserId,
			qusEvent.StockSymbol,
			qusEvent.Price,
			qusEvent.Cryptokey,
			qusEvent.Timestamp)
	case data.ErEvent:
		errEvent := event.ErrEvent
		return data.GetErrorEvent(
			errEvent.Server,
			errEvent.TransactionNum,
			errEvent.Command,
			errEvent.UserId,
			errEvent.StockSymbol,
			errEvent.Funds,
			errEvent.ErrorMessage,
			errEvent.Timestamp)
	default:
		return nil
	}
}

func getXmlEventStringHelper(i interface{}) string {

	var xmlString string
	if xmlstring, err := xml.MarshalIndent(i, "", "    "); err == nil {
		xmlString = string(xmlstring)
		return xmlString
	}
	return xmlString
}

func (service Service) getXMLEventString(event *data.Event) string {

	targetEvent := getEvent(event)
	return getXmlEventStringHelper(targetEvent)
}
