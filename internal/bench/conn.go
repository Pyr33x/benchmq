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
			b.cfg.Client.ClientID = fmt.Sprintf("benchmq-client-%d", id)
			client := mqtt.NewClient(b.cfg)
			defer client.Disconnect()
			log.Printf("Client: %s ― Connected", b.cfg.Client.ClientID)
			if err := client.Connect(); err != nil {
				log.Printf("client: %d failed: %v", id, err)
			}
		}(i)
	}
	b.wg.Wait()
	log.Println(time.Since(now))
}
