package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"time"

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

		go func() {
			time.Sleep(500 * time.Millisecond)
			logger.Info("shutdown", logger.State("completed"))
			os.Exit(0)
		}()

		// Parse flags
		clients, err := cmd.Flags().GetInt("clients")
		if err != nil {
			logger.Error("failed to parse clients flag", logger.ErrorAttr(err))
			return
		}

		delay, err := cmd.Flags().GetInt("delay")
		if err != nil {
			logger.Error("failed to parse delay flag", logger.ErrorAttr(err))
			return
		}

		clean, err := cmd.Flags().GetBool("clean")
		if err != nil {
			logger.Error("failed to parse clean flag", logger.ErrorAttr(err))
			return
		}

		keepalive, err := cmd.Flags().GetUint16("keepalive")
		if err != nil {
			logger.Error("failed to parse keepalive flag", logger.ErrorAttr(err))
			return
		}

		clientID, err := cmd.Flags().GetString("clientID")
		if err != nil {
			logger.Error("failed to parse clientID flag", logger.ErrorAttr(err))
			return
		}

		b, err := bench.NewBenchmark(
			Cfg,
			bench.WithClients(clients),
			bench.WithDelay(delay),
			bench.WithCleanSession(clean),
			bench.WithKeepAlive(keepalive),
			bench.WithClientID(clientID),
		)
		if err != nil {
			logger.Error("Benchmark failed", logger.ErrorAttr(err))
		}

		// Run benchmark
		b.RunConnections()
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
