package config

import (
	"log"

	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
)

func ReadInConfig() *viper.Viper {

	// TODO: add error handling later, return err instead of constructing new one

	v := viper.New()
	v.AddConfigPath("./app/config")
	v.SetConfigType("toml")

	v.SetConfigName("app")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	readInConfigHelper(v, "app.docker")
	readInConfigHelper(v, "app.development")
	readInConfigHelper(v, "app.test")
	readInConfigHelper(v, "messages")
	readInConfigHelper(v, "maps")
	readInConfigHelper(v, "errorCode")

	return v
}

func GetMongoSession(dialInfo *mgo.DialInfo) *mgo.Session {

	c := make(chan *mgo.Session)

	go func() {

		// Connect to MongoDB and establish a connection
		// Only do this once in your application.
		// There is a lot of overhead with this call
		session, err := mgo.DialWithInfo(dialInfo)

		if err != nil {
			log.Fatalln(err)
		} else {
			c <- session
		}

	}()

	mgoSession := <-c

	return mgoSession
}

/*
	Private methods
*/

func readInConfigHelper(v *viper.Viper, fileName string) {

	v.SetConfigName(fileName)
	err := v.MergeInConfig()
	if err != nil {
		log.Fatalln(err)
	}

}
