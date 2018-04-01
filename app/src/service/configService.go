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

func (configService ConfigService) ConnectEventDB(
	session *mgo.Session) *mgo.Database {

	newSession := session.Clone()
	defer newSession.Close()

	dbName := configService.u.GetDBName()
	db := newSession.DB(dbName)

	return db
}

/*
	Private methods
*/

func newConfigService(u *util.Utilities) *ConfigService {

	return &ConfigService{u}
}
