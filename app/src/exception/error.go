package exception

import "fmt"

type ASError struct {
	UUID     UUID
	Msg      string // customized error message
	SysError error
}

func NewASError(uuid UUID, msg string, sysError error) *ASError {
	return &ASError{uuid, msg, sysError}
}

func (e *ASError) Error() *ASError {
	return e
}

func (e *ASError) ErrorMessage() string {
	var formattedMsg string
	if e.SysError == nil {
		formattedMsg = e.Msg
	} else {
		formattedMsg = fmt.Sprintf("%s \n %s", e.Msg, e.SysError.Error())
	}
	return fmt.Sprintf("%s: %s", e.UUID, formattedMsg)
}
