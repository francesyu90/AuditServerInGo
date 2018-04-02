package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../data"
	"../exception"
	"../util"
)

type Repository struct {
	session *mgo.Session
	dbName  string
	u       *util.Utilities
	docName string
	loggers *util.Logger
}

func GetRepository(
	session *mgo.Session,
	u *util.Utilities,
	loggers *util.Logger) *Repository {

	return newRepository(session, u, loggers)
}

func (repo Repository) Insert(event *data.Event) *exception.ASError {

	userID, asErr1 := repo.getUserID(event)
	if asErr1 != nil {
		return asErr1
	}
	event.UserID = userID
	event.ID = bson.NewObjectId()

	repo.loggers.INFO.Println("Event to be inserted: ", event)

	c := make(chan error)

	go func() {

		newSession := repo.session.Clone()
		defer newSession.Close()

		err := newSession.DB(repo.dbName).
			C(repo.docName).
			Insert(&event)

		c <- err
	}()

	err := <-c
	if err != nil {
		asErr := repo.u.GetError(
			exception.AS00008, "db_insert_error", err)
		repo.loggers.ERROR.Println(asErr.ErrorMessage())
		return asErr
	}
	return nil
}

func (repo Repository) UpdateById(event *data.Event) *exception.ASError {

	repo.loggers.INFO.Println("Event to be updated: ", event)

	c := make(chan error)

	go func() {

		newSession := repo.session.Clone()
		defer newSession.Close()

		err := newSession.DB(repo.dbName).
			C(repo.docName).
			UpdateId(event.ID, &event)

		c <- err
	}()

	err := <-c
	if err != nil {
		asErr := repo.u.GetError(
			exception.AS00009, "db_update_by_id_error", err)
		repo.loggers.ERROR.Println(asErr.ErrorMessage())
		return asErr
	}
	return nil
}

func (repo Repository) Delete(event *data.Event) *exception.ASError {

	repo.loggers.INFO.Println("Event to be deleted: ", event)

	c := make(chan error)

	go func() {

		newSession := repo.session.Clone()
		defer newSession.Close()

		err := newSession.DB(repo.dbName).
			C(repo.docName).
			Remove(&event)

		c <- err
	}()

	err := <-c
	if err != nil {
		asErr := repo.u.GetError(
			exception.AS00010, "db_delete_error", err)
		repo.loggers.ERROR.Println(asErr.ErrorMessage())
		return asErr
	}
	return nil
}

func (repo Repository) FindAll() ([]*data.Event, *exception.ASError) {

	var events []*data.Event

	c := make(chan error)

	go func() {

		newSession := repo.session.Clone()
		defer newSession.Close()

		err := newSession.DB(repo.dbName).
			C(repo.docName).
			Find(bson.M{}).
			All(&events)

		c <- err
	}()

	err := <-c
	if err != nil {
		asErr := repo.u.GetError(
			exception.AS00011, "db_find_all_error", err)
		repo.loggers.ERROR.Println(asErr.ErrorMessage())
		return nil, asErr
	}

	repo.loggers.INFO.Println("Events: ", events)
	return events, nil
}

func (repo Repository) FindByUserID(userID string) (
	[]*data.Event, *exception.ASError) {

	var eventsByUser []*data.Event

	c := make(chan error)

	go func() {

		newSession := repo.session.Clone()
		defer newSession.Close()

		err := newSession.DB(repo.dbName).
			C(repo.docName).
			Find(bson.M{"user_id": userID}).
			All(&eventsByUser)

		c <- err
	}()

	err := <-c

	if err != nil {
		asErr := repo.u.GetError(
			exception.AS00012, "db_find_by_user_id_error", err)
		repo.loggers.ERROR.Println(asErr.ErrorMessage())
		return nil, asErr
	}

	repo.loggers.INFO.Println("Events by user: ", eventsByUser)
	return eventsByUser, nil

}

/*
	Private methods
*/
func newRepository(
	session *mgo.Session,
	u *util.Utilities,
	loggers *util.Logger) *Repository {

	dbName := u.GetDBName()
	docName := u.GetDBDocName()
	return &Repository{session, dbName, u, docName, loggers}
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
	repo.loggers.ERROR.Println(asError.ErrorMessage())
	return "", asError
}
