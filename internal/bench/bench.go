package bench

type QoSLevel uint8

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
	QoS0 QoSLevel = 0
	QoS1 QoSLevel = 1
	QoS2 QoSLevel = 2
)

const (
	DefaultDelay        = 100
	DefaultClients      = 100
	DefaultClientID     = "benchmq-client"
	DefaultTopic        = "bench/test"
	DefaultCleanSession = true
	DefaultQoS          = QoS0
	DefaultKeepAlive    = 60
	DefaultHost         = "localhost"
	DefaultPort         = 1883
)

func NewBenchmark(options ...Option) *Bench {
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

	return &bench
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

func WithClientID(clientId string) Option {
	return func(b *Bench) {
		b.ClientID = clientId
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
