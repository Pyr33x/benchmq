package er

import (
	"errors"
	"fmt"
)

var (
	ErrMqttConnectionFailed = errors.New("mqtt connection failed")
	ErrEmptyServerHost      = errors.New("server host cannot be empty")
	ErrInvalidServerPort    = errors.New("server port is invalid")
	ErrUnmarshalFailed      = errors.New("failed to unmarshal config file")
	ErrConfigReadFailed     = errors.New("failed to read config file")
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
