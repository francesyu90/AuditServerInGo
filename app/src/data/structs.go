package data

type AccountTransactionEvent struct {
	Timestamp      int64
	Server         string
	TransactionNum int
	Action         string
	UserId         string
	Funds          string
}

type SystemEvent struct {
	Timestamp      int64
	Server         string
	TransactionNum int
	Command        string
	UserId         string
	StockSymbol    string
	Funds          string
}

type QuoteServerEvent struct {
	Timestamp            int64
	Server               string
	TransactionNum       int
	QuoteServerEventTime int64
	Command              string
	UserId               string
	StockSymbol          string
	Price                string
	Cryptokey            string
}

type ErrorEvent struct {
	Timestamp      int64
	Server         string
	TransactionNum int
	Command        string
	UserId         string
	StockSymbol    string
	Funds          string
	ErrorMessage   string
}
