package config

import (
	"log"

	"github.com/spf13/viper"
)

func ReadInConfig() *viper.Viper {

	// TODO: add error handling later, return err instead of constructing new one

	v := viper.New()
	v.AddConfigPath("./app/config")
	v.SetConfigType("toml")

	v.SetConfigName("app")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	readInConfigHelper(v, "app.docker")
	readInConfigHelper(v, "app.development")
	readInConfigHelper(v, "app.production")
	readInConfigHelper(v, "app.qa")
	readInConfigHelper(v, "messages")
	readInConfigHelper(v, "maps")

	return v
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
