package bench

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/pyr33x/benchmq/internal/mqtt"
	"github.com/pyr33x/benchmq/pkg/logger"
)

func (b *Bench) PublishMessages() {
	start := time.Now()
	b.logger.Info("Started publish benchmark", logger.Int("startTime", int(start.UnixNano())))

	var failed int32
	var succeeded int32

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
			b.logger.Info("Connecting Client", logger.ClientID(id), logger.State("connecting"))

			if err := client.Connect(); err != nil {
				atomic.AddInt32(&failed, int32(b.messageCount))
				b.logger.Error("Client connection failed", logger.ClientID(id), logger.ErrorAttr(err))
				return
			}
			defer client.Disconnect()

			for j := 0; j < b.messageCount; j++ {
				if b.delay > 0 {
					time.Sleep(time.Duration(b.delay) * time.Millisecond)
				}

				err := client.Publish(b.topic, byte(b.qos), b.retained, b.message, func() {
					atomic.AddInt32(&succeeded, 1)
					b.logger.LogPublish(id, b.topic, int(b.qos), b.retained)
				})
				if err != nil {
					atomic.AddInt32(&failed, 1)
					b.logger.Error("Failed to publish message", logger.ErrorAttr(err))
				}
			}
		}(clientID)
	}

	b.wg.Wait()

	elapsed := time.Since(start).Seconds()
	total := b.clients * b.messageCount
	throughput := float64(total) / elapsed

	b.logger.Info("Finished publish benchmark",
		logger.Int("clients", b.clients),
		logger.Int("messagesPerClient", b.messageCount),
		logger.Int("totalMessages", total),
		logger.Int("successful", int(succeeded)),
		logger.Int("failed", int(failed)),
		logger.Float("elapsedSec", elapsed),
		logger.Float("throughputMsgPerSec", throughput),
	)
}
