package bench

import "github.com/pyr33x/benchmq/pkg/er"

type QoSLevel uint8

// Bench represents the benchmark fields
type Bench struct {
	Delay        int
	Clients      int
	ClientID     string
	Topic        string
	CleanSession bool
	QoS          QoSLevel
	KeepAlive    uint16
	Host         string
	Port         uint16
}

type Option func(*Bench)

const (
	QoS0 QoSLevel = 0 // QoS At Most Once
	QoS1 QoSLevel = 1 // QoS At Lease Once
	QoS2 QoSLevel = 2 // QoS Exactly Once
)

const (
	DefaultDelay        = 100              // Default delay between connection
	DefaultClients      = 100              // Default clients to connect
	DefaultClientID     = "benchmq-client" // Default client id
	DefaultTopic        = "bench/test"     // Default publish/subscribe topic
	DefaultCleanSession = true             // Default clean session state
	DefaultQoS          = QoS0             // Default QoS level
	DefaultKeepAlive    = 60               // Default connection keep alive
	DefaultHost         = "localhost"      // Default broker host
	DefaultPort         = 1883             // Default broker port
)

// NewBenchmark constructor initializes the bench struct
func NewBenchmark(options ...Option) (*Bench, error) {
	bench := Bench{
		Delay:        DefaultDelay,
		Clients:      DefaultClients,
		ClientID:     DefaultClientID,
		Topic:        DefaultTopic,
		CleanSession: DefaultCleanSession,
		QoS:          DefaultQoS,
		KeepAlive:    DefaultKeepAlive,
		Host:         DefaultHost,
		Port:         DefaultPort,
	}

	for _, option := range options {
		if option != nil {
			option(&bench)
		}
	}

	if err := bench.Validate(); err != nil {
		return nil, err
	}

	return &bench, nil
}

// Validate checks semantic correctness of the benchmark configuration
func (b Bench) Validate() error {
	if b.Clients <= 0 {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrInvalidClients,
			Raw:     er.ErrInvalidClients,
		}
	}
	if b.Delay < 0 {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrInvalidDelay,
			Raw:     er.ErrInvalidDelay,
		}
	}
	if b.Host == "" {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrEmptyHost,
			Raw:     er.ErrEmptyHost,
		}
	}
	if b.Topic == "" {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrEmptyTopic,
			Raw:     er.ErrEmptyTopic,
		}
	}
	if b.Port == 0 {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrInvalidPort,
			Raw:     er.ErrInvalidPort,
		}
	}
	if b.QoS > QoS2 {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrInvalidQoS,
			Raw:     er.ErrInvalidQoS,
		}
	}
	return nil
}

func WithDelay(delay int) Option {
	return func(b *Bench) {
		b.Delay = delay
	}
}

func WithClients(clients int) Option {
	return func(b *Bench) {
		b.Clients = clients
	}
}

func WithClientID(clientID string) Option {
	return func(b *Bench) {
		b.ClientID = clientID
	}
}

func WithTopic(topic string) Option {
	return func(b *Bench) {
		b.Topic = topic
	}
}

func WithCleanSession(cleanSession bool) Option {
	return func(b *Bench) {
		b.CleanSession = cleanSession
	}
}

func WithQoS(qos uint16) Option {
	return func(b *Bench) {
		b.QoS = QoSLevel(qos)
	}
}

func WithKeepAlive(keepAlive uint16) Option {
	return func(b *Bench) {
		b.KeepAlive = keepAlive
	}
}

func WithHost(host string) Option {
	return func(b *Bench) {
		b.Host = host
	}
}

func WithPort(port uint16) Option {
	return func(b *Bench) {
		b.Port = port
	}
}
