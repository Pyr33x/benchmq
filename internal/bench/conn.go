package bench

import (
	"fmt"
	"log"
	"time"

	"github.com/pyr33x/benchmq/internal/mqtt"
)

func (b *Bench) RunConnections() {
	now := time.Now()

	for i := 0; i < b.clients; i++ {
		b.wg.Add(1)
		go func(id int) {
			defer b.wg.Done()

			cfg := *b.cfg
			cfg.Client.ClientID = fmt.Sprintf("benchmq-client-%d", id)
			client := mqtt.NewClient(&cfg)
			defer client.Disconnect()

			log.Printf("Client: %s ― Connecting", cfg.Client.ClientID)
			if err := client.Connect(); err != nil {
				log.Printf("Client: %d ― Failed: %v", id, err)
				return
			}
			log.Printf("Client: %s ― Connected", cfg.Client.ClientID)

		}(i)
	}

	b.wg.Wait()
	log.Println(time.Since(now))
}
