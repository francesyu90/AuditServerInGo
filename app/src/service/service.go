package service

import (
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"

	"../util"
)

type Service struct {
	u *util.Utilities
}

func GetService(u *util.Utilities) *Service {
	return newService(u)
}

func (service Service) ProcessReqBody(
	c *gin.Context, i interface{}) <-chan error {

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

/*
	Private methods
*/

func newService(u *util.Utilities) *Service {
	return &Service{u}
}
