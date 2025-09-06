package er

import (
	"errors"
	"fmt"
)

var (
	ErrMqttConnectionFailed = errors.New("mqtt connection failed")
)

type Error struct {
	Context string
	Message error
}

func (e *Error) Error() string {
	return fmt.Sprintf("context: %s, error: %v", e.Context, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Message
}
