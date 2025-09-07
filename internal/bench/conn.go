package bench

import (
	"fmt"
	"time"

	"github.com/pyr33x/benchmq/internal/mqtt"
	"github.com/pyr33x/benchmq/pkg/logger"
)

func (b *Bench) RunConnections() {
	start := time.Now()
	b.logger.Info("Started connection benchmark", logger.Int("current_time", int(start.UnixNano())))

	for i := 0; i < b.clients; i++ {
		b.wg.Add(1)
		go func(id int) {
			defer b.wg.Done()

			cfg := *b.cfg
			cfg.Client.ClientID = fmt.Sprintf("benchmq-client-%d", id)
			client := mqtt.NewClient(&cfg)
			defer client.Disconnect()

			b.logger.Info("Connecting Client", logger.ClientID(cfg.Client.ClientID), logger.State("connecting"))
			if err := client.Connect(); err != nil {
				b.logger.Error("Couldn't establish client", logger.ClientID(cfg.Client.ClientID), logger.State("failed"))
				return
			}
			b.logger.LogClientConnection(cfg.Client.ClientID)

		}(i)
	}

	b.wg.Wait()
	b.logger.Info("Finished connection benchmark", logger.TrackTime(start))
}
