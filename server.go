package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"

	"./app/src/config"
	"./app/src/controller"
	"./app/src/util"
)

func getMainEngine(u *util.Utilities) (*gin.Engine, string) {
	router := gin.Default()
	controller := controller.GetController(u)
	api := router.Group(u.GetStringConfigValue("uri.api"))
	{
		api.POST(u.GetStringConfigValue("test.uri.testing"), controller.Testing)
	}

	portStr := fmt.Sprintf(":%d", u.GetIntConfigValue("environment.port"))

	return router, portStr
}

func setUpHelper() *util.Utilities {

	v := config.ReadInConfig()
	return util.NewUtilities(v)
}

func setUp() {

	u := setUpHelper()
	router, portStr := getMainEngine(u)
	router.Run(portStr)
}

func test() {
	u := setUpHelper()

	dialInfo := u.GetDBDialInfo()
	// url := u.GetDBUrl()

	// dialInfo := &mgo.DialInfo{
	// 	Addrs:   []string{url},
	// 	Timeout: 10 * time.Second,
	// }

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(session)
}

func main() {

	setUp()

	// test()
}
