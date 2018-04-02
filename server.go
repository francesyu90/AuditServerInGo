package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"./app/src/config"
	"./app/src/controller"
	"./app/src/util"
)

func getMainEngine(u *util.Utilities, loggers *util.Logger) (
	*gin.Engine, string) {

	router := gin.Default()

	controller := controller.GetController(u, loggers)

	api := router.Group(u.GetStringConfigValue("uri.api"))
	{
		api.GET(
			u.GetStringConfigValue("test.uri.testing"),
			controller.Testing)
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
	loggers, asErr := util.InitLoggers(u)
	if asErr != nil {
		log.Fatalln(asErr)
	}
	router, portStr := getMainEngine(u, loggers)
	router.Run(portStr)
}

func main() {

	setUp()

}
