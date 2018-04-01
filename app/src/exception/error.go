package exception

import "fmt"

type ASError struct {
	UUID   UUID
	ErrMsg string // system error
	Msg    string // customized error message
}

func NewASError(uuid UUID, errMsg string, msg string) *ASError {
	return &ASError{uuid, errMsg, msg}
}

func (e *ASError) Error() *ASError {
	return e
}

func (e *ASError) ErrorMessage() string {
	var formattedMsg string
	if e.ErrMsg == "" {
		formattedMsg = e.Msg
	} else if e.Msg == "" {
		formattedMsg = e.ErrMsg
	} else {
		formattedMsg = fmt.Sprintf("%s \n %s", e.Msg, e.ErrMsg)
	}
	return fmt.Sprintf("%s: %s", e.UUID, formattedMsg)
}
