package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "benchmq",
	Short: "BenchMQ is a simple, fast, and lightweight CLI to benchmark your MQTT broker with ease.",
	Long:  `BenchMQ is a simple, fast, and open-source CLI tool for benchmarking MQTT brokers. Measure throughput, latency, and stability of your MQTT setup with ease.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
