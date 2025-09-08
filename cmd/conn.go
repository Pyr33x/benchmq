package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pyr33x/benchmq/internal/bench"
	"github.com/pyr33x/benchmq/pkg/logger"
	"github.com/spf13/cobra"
)

var connCmd = &cobra.Command{
	Use:   "conn",
	Short: "Run a connection benchmark against the configured MQTT broker.",
	Long:  `Opens N concurrent MQTT connections (from config or flags) to measure connection throughput, failures, and timing.`,
	Run: func(cmd *cobra.Command, args []string) {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		// Parse flags
		clients, _ := cmd.Flags().GetInt("clients")
		delay, _ := cmd.Flags().GetInt("delay")
		clean, _ := cmd.Flags().GetBool("clean")
		keepalive, _ := cmd.Flags().GetUint16("keepalive")
		clientID, _ := cmd.Flags().GetString("clientID")

		// Create benchmark
		b, err := bench.NewBenchmark(
			Cfg,
			bench.WithClients(clients),
			bench.WithDelay(delay),
			bench.WithCleanSession(clean),
			bench.WithKeepAlive(keepalive),
			bench.WithClientID(clientID),
		)
		if err != nil {
			logger.Error("Failed to create benchmark", logger.ErrorAttr(err))
			return
		}

		// Run benchmark in a goroutine so we can wait for shutdown
		done := make(chan struct{})
		go func() {
			b.RunConnections()
			close(done)
		}()

		select {
		case <-sigs:
			logger.Info("Received shutdown signal", logger.State("terminated"))
			return
		case <-done:
			logger.Info("Connection benchmark completed", logger.State("completed"))
		}
	},
}

func init() {
	rootCmd.AddCommand(connCmd)

	// Register flags
	connCmd.Flags().StringP("clientID", "i", "benchmq-client", "Client ID for MQTT connections")
	connCmd.Flags().IntP("clients", "c", 100, "Number of concurrent clients to connect")
	connCmd.Flags().IntP("delay", "d", 1000, "Delay between each client connection in milliseconds")
	connCmd.Flags().BoolP("clean", "x", true, "Clean previous session when connecting")
	connCmd.Flags().Uint16P("keepalive", "k", 60, "Keepalive interval in seconds")
}
