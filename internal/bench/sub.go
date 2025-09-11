package bench

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/pyr33x/benchmq/internal/mqtt"
	"github.com/pyr33x/benchmq/pkg/logger"
)

func (b *Bench) Subscribe() {
	start := time.Now()
	b.logger.Info("Started subscribe benchmark", logger.String("start", start.Format(time.RFC3339Nano)))

	var received int64
	var failed int64

	for i := 0; i < b.clients; i++ {
		b.wg.Add(1)

		clientID := fmt.Sprintf("%s-%d", b.clientID, i)
		go func(id string) {
			defer b.wg.Done()

			cfg := *b.cfg
			cfg.Client.ClientID = id
			cfg.Client.CleanSession = *b.cleanSession
			cfg.Client.KeepAlive = b.keepAlive
			client := mqtt.NewClient(&cfg)

			b.logger.Info("Connecting subscriber", logger.ClientID(id), logger.State("connecting"))
			if err := client.Connect(); err != nil {
				atomic.AddInt64(&failed, 1)
				b.logger.Error("Subscriber connection failed", logger.ClientID(id), logger.ErrorAttr(err))
				return
			}
			defer client.Disconnect()

			err := client.Subscribe(b.topic, byte(b.qos), b.retained, func(payload string) {
				atomic.AddInt64(&received, 1)
				b.logger.LogSubscribe(id, b.topic, int(b.qos), logger.String("payload", payload))
			})
			if err != nil {
				atomic.AddInt64(&failed, 1)
				b.logger.Error("Failed to subscribe", logger.ClientID(id), logger.ErrorAttr(err))
				return
			}

			if b.delay > 0 {
				time.Sleep(time.Duration(b.delay) * time.Millisecond * time.Duration(b.messageCount))
			} else {
				time.Sleep(time.Second * 5)
			}
		}(clientID)
	}

	b.wg.Wait()

	elapsed := time.Since(start).Seconds()
	expected := int64(b.clients) * int64(b.messageCount)
	throughput := float64(received) / elapsed
	b.logger.Info("Finished subscribe benchmark",
		logger.Int("clients", b.clients),
		logger.Any("expectedMessages", expected),
		logger.Any("received", received),
		logger.Any("failed", failed),
		logger.Float("elapsedSec", elapsed),
		logger.Float("throughputMsgPerSec", throughput),
	)
}
