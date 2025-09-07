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

func NewBenchmark(options ...Option) *Bench {
	bench := Bench{
		Delay:        100,
		Clients:      100,
		ClientID:     "benchmq-client",
		Topic:        "bench/test",
		CleanSession: true,
		QoS:          0,
		KeepAlive:    60,
		Host:         "localhost",
		Port:         1883,
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
