package mqtt

import (
	"fmt"
	"sync"
	"time"

	mq "github.com/eclipse/paho.mqtt.golang"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pyr33x/benchmq/pkg/config"
	"github.com/pyr33x/benchmq/pkg/er"
	"github.com/pyr33x/benchmq/pkg/logger"
)

// Adapter represents an MQTT adapter instance
type Adapter struct {
	client mq.Client
	wg     sync.WaitGroup
}

// NewClient creates a new MQTT adapter instance
func NewClient(cfg *config.Config) *Adapter {
	// Initialize MQTT client options
	opts := mq.NewClientOptions()

	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.Server.Host, cfg.Server.Port))
	opts.SetClientID(cfg.Client.ClientID)
	opts.SetKeepAlive(time.Duration(cfg.Client.KeepAlive) * time.Second)
	opts.SetCleanSession(cfg.Client.CleanSession)
	opts.SetUsername(cfg.Client.Username)
	opts.SetPassword(cfg.Client.Password)
	opts.SetProtocolVersion(4) // Default set to MQTT 3.1.1

	// Create a new MQTT client instance
	client := mq.NewClient(opts)

	// Return the initialized MQTT adapter
	return &Adapter{client: client}
}

// Connect establishes a connection to the MQTT broker
func (a *Adapter) Connect() error {
	if token := a.client.Connect(); token.Wait() && token.Error() != nil {
		tErr := token.Error()
		return &er.Error{
			Package: "MQTT",
			Func:    "Connect",
			Message: er.ErrMqttConnectionFailed,
			Raw:     tErr,
		}
	}
	return nil
}

// Publish publishes a message to the specified topic with the given QoS level and retention flag
func (a *Adapter) Publish(topic string, qos byte, retained bool, payload any, callback func()) error {
	if callback == nil {
		return &er.Error{
			Package: "MQTT",
			Func:    "Publish",
			Message: er.ErrNilCallback,
		}
	}

	if err := a.Validate(topic, qos); err != nil {
		return err
	}

	token := a.client.Publish(topic, qos, retained, payload)
	token.Wait()

	if err := token.Error(); err != nil {
		return &er.Error{
			Package: "MQTT",
			Func:    "Publish",
			Message: er.ErrPublishFailed,
			Raw:     err,
		}
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		callback()
	}()

	return nil
}

// Unsubscribe unsubscribes from the specified topic
func (a *Adapter) Unsubscribe(topic string) error {
	if err := a.Validate(topic, 0); err != nil {
		return err
	}

	token := a.client.Unsubscribe(topic)
	token.Wait()

	if err := token.Error(); err != nil {
		return &er.Error{
			Package: "MQTT",
			Func:    "Unsubscribe",
			Message: er.ErrUnsubscribeFailed,
			Raw:     err,
		}
	}

	return nil
}

// Subscribe subscribes to the specified topic with the given QoS level and retention flag
func (a *Adapter) Subscribe(topic string, qos byte, retained bool, callback func(payload string)) error {
	if callback == nil {
		return &er.Error{
			Package: "MQTT",
			Func:    "Subscribe",
			Message: er.ErrNilCallback,
		}
	}

	if err := a.Validate(topic, qos); err != nil {
		return err
	}

	token := a.client.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logger.Error("panic in subscription callback",
						logger.Any("recover", r),
						logger.String("topic", topic),
					)
				}
			}()
			callback(payload)
		}()
	})
	token.Wait()

	if err := token.Error(); err != nil {
		return &er.Error{
			Package: "MQTT",
			Func:    "Subscribe",
			Message: er.ErrSubscribeFailed,
			Raw:     err,
		}
	}

	return nil
}

// Validate validates the topic and QoS level
func (a *Adapter) Validate(topic string, qos byte) error {
	if topic == "" {
		return &er.Error{
			Package: "MQTT",
			Func:    "Validate",
			Message: er.ErrEmptyTopic,
		}
	}
	if qos > 2 {
		return &er.Error{
			Package: "MQTT",
			Func:    "Validate",
			Message: er.ErrInvalidQoS,
		}
	}

	return nil
}

// Disconnect disconnects the client from the MQTT broker
func (a *Adapter) Disconnect() {
	a.client.Disconnect(200)
	a.wg.Wait()
}
