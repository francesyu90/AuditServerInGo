package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"

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

func getMainEngine() (*gin.Engine, string) {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/testing", testing)
	}

	portStr := fmt.Sprintf(":%d", 8082)

	return router, portStr
}

func setUp() {
	router, portStr := getMainEngine()
	router.Run(portStr)
}

func main() {
	setUp()
}
