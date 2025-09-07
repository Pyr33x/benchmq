package er

import (
	"errors"
	"fmt"
)

var (
	ErrMqttConnectionFailed = errors.New("mqtt connection failed")
)

type Error struct {
	Package string
	Func    string
	Message error
}

func (e *Error) Error() string {
	return fmt.Sprintf("package: %s, func: %s, error: %v", e.Package, e.Func, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Message
}
