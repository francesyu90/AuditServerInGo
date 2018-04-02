package data

import (
	"encoding/xml"
	"fmt"
)

/*
	XML structs
*/

// UserCommand infomration necessary for logging user command
type UCommand struct {
	XMLName        xml.Name `xml:"userCommand"`
	Timestamp      string   `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum string   `xml:"transactionNum"`
	Command        string   `xml:"command"`
	Username       string   `xml:"username"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
}

// AccountTransaction infomration necessary for logging account transaction
type ATransaction struct {
	XMLName        xml.Name `xml:"accountTransaction"`
	Timestamp      string   `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum string   `xml:"transactionNum"`
	Action         string   `xml:"action"`
	Username       string   `xml:"username"`
	Funds          string   `xml:"funds"`
}

// SystemEvent information necessary for logging system event
type SEvent struct {
	XMLName        xml.Name `xml:"systemEvent"`
	Timestamp      string   `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum string   `xml:"transactionNum"`
	Command        string   `xml:"command,omitempty"`
	Username       string   `xml:"username"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
}

// QuoteServer information necessary for logging quote server hit
type QServer struct {
	XMLName         xml.Name `xml:"quoteServer"`
	Timestamp       string   `xml:"timestamp"`
	Server          string   `xml:"server"`
	TransactionNum  string   `xml:"transactionNum"`
	QuoteServerTime string   `xml:"quoteServerTime"`
	Command         string   `xml:"command,omitempty"`
	Username        string   `xml:"username"`
	StockSymbol     string   `xml:"stockSymbol"`
	Price           string   `xml:"price"`
	Cryptokey       string   `xml:"cryptokey"`
}

type EEvent struct {
	XMLName        xml.Name `xml:"errorEvent"`
	Timestamp      string   `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum string   `xml:"transactionNum"`
	Command        string   `xml:"command,omitempty"`
	Username       string   `xml:"username"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
	ErrorMessage   string   `xml:"errorMessage"`
}

func GetUserCommand(
	server string,
	transactionNum int,
	command string,
	username string,
	stockSymbol string,
	funds float64,
	timestamp int64) UCommand {

	fundsAsString := getFundsAsString(funds)

	return UCommand{
		Timestamp:      fmt.Sprint(timestamp),
		Server:         server,
		TransactionNum: fmt.Sprint(transactionNum),
		Command:        command,
		Username:       username,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString}
}

func GetAccountTransaction(
	server string,
	transactionNum int,
	action string,
	username string,
	funds string,
	timestamp int64) ATransaction {

	return ATransaction{
		Timestamp:      fmt.Sprint(timestamp),
		Server:         server,
		TransactionNum: fmt.Sprint(transactionNum),
		Action:         action,
		Username:       username,
		Funds:          funds}
}

func GetSystemEvent(
	server string,
	transactionNum int,
	command string,
	username string,
	stockSymbol string,
	funds string,
	timestamp int64) SEvent {

	return SEvent{
		Timestamp:      fmt.Sprint(timestamp),
		Server:         server,
		TransactionNum: fmt.Sprint(transactionNum),
		Command:        command,
		Username:       username,
		StockSymbol:    stockSymbol,
		Funds:          funds}
}

func GetQuoteServer(
	server string,
	transactionNum int,
	quoteServerTime int64,
	command string,
	username string,
	stockSymbol string,
	price string,
	cryptokey string,
	timestamp int64) QServer {

	return QServer{
		Timestamp:       fmt.Sprint(timestamp),
		Server:          server,
		TransactionNum:  fmt.Sprint(transactionNum),
		QuoteServerTime: fmt.Sprint(quoteServerTime),
		Command:         command,
		Username:        username,
		StockSymbol:     stockSymbol,
		Price:           price,
		Cryptokey:       cryptokey}
}

func GetErrorEvent(
	server string,
	transactionNum int,
	command string,
	username string,
	stockSymbol string,
	funds string,
	errorMessage string,
	timestamp int64) EEvent {

	return EEvent{
		Timestamp:      fmt.Sprint(timestamp),
		Server:         server,
		TransactionNum: fmt.Sprint(transactionNum),
		Command:        command,
		Username:       username,
		StockSymbol:    stockSymbol,
		Funds:          funds,
		ErrorMessage:   errorMessage}
}

/*
	Private methods
*/
func getFundsAsString(amount float64) string {
	if amount == 0 {
		return ""
	}
	return fmt.Sprintf("%.2f", float64(amount))
}
