package bench

import (
	"sync"

	"github.com/pyr33x/benchmq/pkg/config"
	"github.com/pyr33x/benchmq/pkg/er"
)

type QoSLevel uint8

// Bench represents the benchmark fields
type Bench struct {
	delay        int
	clients      int
	clientID     string
	topic        string
	cleanSession bool
	qos          QoSLevel
	keepAlive    uint16
	host         string
	port         uint16
	wg           sync.WaitGroup // Wait Group
	cfg          *config.Config // Config
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
func NewBenchmark(cfg *config.Config, options ...Option) (*Bench, error) {
	bench := Bench{
		delay:        DefaultDelay,
		clients:      DefaultClients,
		clientID:     DefaultClientID,
		topic:        DefaultTopic,
		cleanSession: DefaultCleanSession,
		qos:          DefaultQoS,
		keepAlive:    DefaultKeepAlive,
		host:         DefaultHost,
		port:         DefaultPort,
		cfg:          cfg,
	}

	for _, option := range options {
		if option != nil {
			option(&bench)
		}
	}

	if err := bench.validate(); err != nil {
		return nil, err
	}

	return &bench, nil
}

// Validate checks semantic correctness of the benchmark configuration
func (b *Bench) validate() error {
	if b.clients <= 0 {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrInvalidClients,
			Raw:     er.ErrInvalidClients,
		}
	}
	if b.delay < 0 {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrInvalidDelay,
			Raw:     er.ErrInvalidDelay,
		}
	}
	if b.host == "" {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrEmptyHost,
			Raw:     er.ErrEmptyHost,
		}
	}
	if b.topic == "" {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrEmptyTopic,
			Raw:     er.ErrEmptyTopic,
		}
	}
	if b.port == 0 {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrInvalidPort,
			Raw:     er.ErrInvalidPort,
		}
	}
	if b.qos > QoS2 {
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
		b.delay = delay
	}
}

func WithClients(clients int) Option {
	return func(b *Bench) {
		b.clients = clients
	}
}

func WithClientID(clientID string) Option {
	return func(b *Bench) {
		b.clientID = clientID
	}
}

func WithTopic(topic string) Option {
	return func(b *Bench) {
		b.topic = topic
	}
}

func WithCleanSession(cleanSession bool) Option {
	return func(b *Bench) {
		b.cleanSession = cleanSession
	}
}

func WithQoS(qos uint16) Option {
	return func(b *Bench) {
		b.qos = QoSLevel(qos)
	}
}

func WithKeepAlive(keepAlive uint16) Option {
	return func(b *Bench) {
		b.keepAlive = keepAlive
	}
}

func WithHost(host string) Option {
	return func(b *Bench) {
		b.host = host
	}
}

func WithPort(port uint16) Option {
	return func(b *Bench) {
		b.port = port
	}
}
