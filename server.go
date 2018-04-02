package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"./app/src/config"
	"./app/src/controller"
	"./app/src/exception"
	"./app/src/util"
)

func getMainEngine(u *util.Utilities) (*gin.Engine, string) {
	router := gin.Default()
	controller := controller.GetController(u)
	api := router.Group(u.GetStringConfigValue("uri.api"))
	{
		api.GET(u.GetStringConfigValue("test.uri.testing"), controller.Testing)
		api.POST(
			u.GetStringConfigValue("uri.quote_server"),
			controller.HandleQSEvent)
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
	u.GetError(exception.AS00001, "db_conn_err", nil)
}

func main() {

	setUp()

	// test()

}
