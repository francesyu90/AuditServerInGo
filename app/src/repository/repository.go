package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../data"
	"../exception"
	"../util"
)

type Repository struct {
	session  *mgo.Session
	database *mgo.Database
	u        *util.Utilities
	docName  string
}

func GetRepository(session *mgo.Session, u *util.Utilities) *Repository {
	return newRepository(session, u)
}

func (repo Repository) Insert(event *data.Event) *exception.ASError {

	userID, asErr1 := repo.getUserID(event)
	if asErr1 != nil {
		return asErr1
	}
	event.UserID = userID
	event.ID = bson.NewObjectId()

	c := make(chan error)

	go func() {
		err := repo.database.C(repo.docName).Insert(&event)
		c <- err
	}()

	err := <-c
	if err != nil {
		asErr := repo.u.GetError(
			exception.AS00008, "db_insert_error", err)
		return asErr
	}
	return nil
}

func (repo Repository) UpdateById(event *data.Event) *exception.ASError {

	c := make(chan error)

	go func() {
		err :=
			repo.database.C(repo.docName).UpdateId(event.ID, &event)
		c <- err
	}()

	err := <-c
	if err != nil {
		asErr := repo.u.GetError(
			exception.AS00009, "db_update_by_id_error", err)
		return asErr
	}
	return nil
}

func (repo Repository) Delete(event *data.Event) *exception.ASError {

	c := make(chan error)

	go func() {
		err := repo.database.C(repo.docName).Remove(&event)
		c <- err
	}()

	err := <-c
	if err != nil {
		asErr := repo.u.GetError(
			exception.AS00010, "db_delete_error", err)
		return asErr
	}
	return nil
}

func (repo Repository) FindAll() ([]*data.Event, *exception.ASError) {

	var events []*data.Event

	c := make(chan error)

	go func() {
		err := repo.database.C(repo.docName).Find(bson.M{}).All(&events)
		c <- err
	}()

	err := <-c
	if err != nil {
		asErr := repo.u.GetError(
			exception.AS00011, "db_find_all_error", err)
		return nil, asErr
	}
	return events, nil
}

func (repo Repository) FindByUserID(userID string) (
	[]*data.Event, *exception.ASError) {

	var eventsByUser []*data.Event

	c := make(chan error)

	go func() {

		err := repo.database.
			C(repo.docName).
			Find(bson.M{"user_id": userID}).
			All(&eventsByUser)

		c <- err
	}()

	err := <-c

	if err != nil {
		asErr := repo.u.GetError(
			exception.AS00012, "db_find_by_user_id_error", err)
		return nil, asErr
	}
	return eventsByUser, nil

}

/*
	Private methods
*/
func newRepository(session *mgo.Session, u *util.Utilities) *Repository {

	newSession := session.Clone()

	dbName := u.GetDBName()
	db := newSession.DB(dbName)
	docName := u.GetDBDocName()
	return &Repository{session, db, u, docName}
}

func (repo Repository) getUserID(event *data.Event) (string, *exception.ASError) {

	if event.AcctTxnEvent != nil {
		return event.AcctTxnEvent.UserId, nil
	} else if event.SysEvent != nil {
		return event.SysEvent.UserId, nil
	} else if event.QsEvent != nil {
		return event.QsEvent.UserId, nil
	} else if event.ErrEvent != nil {
		return event.ErrEvent.UserId, nil
	}
	asError := repo.u.GetError(
		exception.AS00013, "no_user_id_available_error", nil)
	return "", asError
}
