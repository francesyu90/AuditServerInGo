package service

import (
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
