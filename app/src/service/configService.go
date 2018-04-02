package service

import (
	"gopkg.in/mgo.v2"

	"../config"
	"../exception"
	"../util"
)

type ConfigService struct {
	u *util.Utilities
}

func GetConfigService(u *util.Utilities) *ConfigService {
	return newConfigService(u)
}

func (configService ConfigService) GetDBConn() (
	*mgo.Session, *exception.ASError) {

	dialInfo, asError := configService.u.GetDBDialInfo()
	if asError != nil {
		return nil, asError
	}
	return config.GetMongoSession(dialInfo), nil
}

/*
	Private methods
*/

func newConfigService(u *util.Utilities) *ConfigService {

	return &ConfigService{u}
}
