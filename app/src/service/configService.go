package service

import (
	"gopkg.in/mgo.v2"

	"../config"
	"../util"
)

type ConfigService struct {
	u *util.Utilities
}

func GetConfigService(u *util.Utilities) *ConfigService {
	return newConfigService(u)
}

func (configService ConfigService) GetDBConn() *mgo.Session {

	dialInfo := configService.u.GetDBDialInfo()
	return config.GetMongoSession(dialInfo)
}

/*
	Private methods
*/

func newConfigService(u *util.Utilities) *ConfigService {

	return &ConfigService{u}
}
