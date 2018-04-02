package service

import (
	"gopkg.in/mgo.v2"

	"../config"
	"../exception"
	"../util"
)

type ConfigService struct {
	u       *util.Utilities
	loggers *util.Logger
}

func GetConfigService(
	u *util.Utilities, loggers *util.Logger) *ConfigService {
	return newConfigService(u, loggers)
}

func (configService ConfigService) GetDBConn() (
	*mgo.Session, *exception.ASError) {

	dialInfo, asError := configService.u.GetDBDialInfo()
	if asError != nil {
		configService.loggers.ERROR.Println(asError.ErrorMessage())
		return nil, asError
	}
	configService.loggers.INFO.Println("Dial Info: ", dialInfo)
	return config.GetMongoSession(dialInfo), nil
}

/*
	Private methods
*/

func newConfigService(
	u *util.Utilities, loggers *util.Logger) *ConfigService {

	return &ConfigService{u, loggers}
}
