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
	ErrInvalidQoS           = errors.New("bench: invalid QoS (must be 0, 1, or 2)")
	ErrInvalidClients       = errors.New("bench: clients must be > 0")
	ErrInvalidDelay         = errors.New("bench: delay must be >= 0")
	ErrInvalidPort          = errors.New("bench: port must be in 1..65535")
	ErrEmptyHost            = errors.New("bench: host must be non-empty")
	ErrEmptyTopic           = errors.New("bench: topic must be non-empty")
	ErrNilConfig            = errors.New("bench: config cannot be nil")
	ErrPublishFailed        = errors.New("mqtt: failed to publish")
)

type Error struct {
	Package string
	Func    string
	Message error
	Raw     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("package: %s, func: %s, error: %v, rawError: %v", e.Package, e.Func, e.Message, e.Raw)
}

func (e *Error) Unwrap() error {
	return e.Message
}
