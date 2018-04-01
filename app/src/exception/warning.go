package exception

import "fmt"

type ASWarning struct {
	UUID UUID
	Msg  string // customized error message
}

func NewASWarning(uuid UUID, msg string) *ASWarning {
	return &ASWarning{uuid, msg}
}

func (w *ASWarning) Warning() *ASWarning {
	return w
}

func (w *ASWarning) WarningMessage() string {
	return fmt.Sprintf("%s: %s", w.UUID, w.Msg)
}
