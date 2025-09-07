package cmd

import (
	"os"

	"github.com/pyr33x/benchmq/pkg/config"
	"github.com/pyr33x/benchmq/pkg/logger"
	"github.com/spf13/cobra"
)

var Cfg *config.Config

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

func init() {
	Cfg, _ = config.InitializeCfg()
	switch Cfg.Environment {
	case "production":
		logger.InitGlobalLogger(logger.ProductionConfig())
	case "development":
		logger.InitGlobalLogger(logger.DevelopmentConfig())
	default:
		logger.InitGlobalLogger(logger.DevelopmentConfig())
		logger.Warn("Invalid server environment config value, assigning default development.")
	}
}
