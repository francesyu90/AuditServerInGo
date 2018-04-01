package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"./app/src/config"
	"./app/src/data"
	"./app/src/util"
)

func testing(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// TODO
	}

	fmt.Println(string(body))

	var acctTxnEvent data.QuoteServerEvent

	util.UnserializeObject(body, &acctTxnEvent)

	c.JSON(200, acctTxnEvent)

}

func getMainEngine(u *util.Utilities) (*gin.Engine, string) {
	router := gin.Default()
	api := router.Group(u.GetStringConfigValue("uri.api"))
	{
		api.POST(u.GetStringConfigValue("test.uri.testing"), testing)
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

func main() {
	setUp()
}