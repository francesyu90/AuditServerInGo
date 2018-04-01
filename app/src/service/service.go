package service

import (
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"

	"../exception"
	"../util"
)

type Service struct {
	u *util.Utilities
}

func GetService(u *util.Utilities) *Service {
	return newService(u)
}

func (service Service) ProcessReqBody(
	c *gin.Context, i interface{}) <-chan *exception.ASError {

	ch := make(chan *exception.ASError)

	go func() {

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			asErr := service.u.GetError(exception.AS00003, "", err)
			ch <- asErr
		}

		log.Println(string(body))

		asErr1 := service.u.UnserializeObject(body, &i)
		if asErr1 != nil {
			ch <- asErr1
		} else {
			ch <- nil
		}

	}()

	return ch
}

/*
	Private methods
*/

func newService(u *util.Utilities) *Service {
	return &Service{u}
}
