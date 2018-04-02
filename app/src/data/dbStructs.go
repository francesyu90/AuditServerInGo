package data

import "gopkg.in/mgo.v2/bson"

type EventType int

const (
	AcTxnEvent = iota
	SyEvent
	QuSEvent
	ErEvent
)

type Event struct {
	ID           bson.ObjectId            `bson:"_id"`
	UserID       string                   `bson:"user_id"`
	EventType    EventType                `bson:"event_type"`
	AcctTxnEvent *AccountTransactionEvent `bson:"account_transaction_event"`
	SysEvent     *SystemEvent             `bson:"system_event"`
	QsEvent      *QuoteServerEvent        `bson:"quote_server_event"`
	ErrEvent     *ErrorEvent              `bson:"error_event"`
}
