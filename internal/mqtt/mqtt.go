package mqtt

import (
	"fmt"
	"sync"
	"time"

	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/pyr33x/benchmq/pkg/config"
	"github.com/pyr33x/benchmq/pkg/er"
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
	token := a.client.Publish(topic, qos, retained, payload)
	a.wg.Add(1)

	go func() {
		defer a.wg.Done()
		callback()
	}()
	token.Wait()

	if token.Error() != nil {
		return &er.Error{
			Package: "MQTT",
			Func:    "Publish",
			Message: er.ErrPublishFailed,
			Raw:     token.Error(),
		}
	}

	return nil
}

// Disconnect disconnects the client from the MQTT broker
func (a *Adapter) Disconnect() {
	a.client.Disconnect(200)
	a.wg.Wait()
}
