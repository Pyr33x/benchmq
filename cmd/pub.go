/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pyr33x/benchmq/internal/bench"
	"github.com/pyr33x/benchmq/pkg/logger"
	"github.com/spf13/cobra"
)

// pubCmd represents the pub command
var pubCmd = &cobra.Command{
	Use:   "pub",
	Short: "Publish messages to a topic with specified parameters",
	Long: `Publish messages to a topic with specified parameters.

Parameters:
    - clients: Number of concurrent clients
    - count: Number of messages to publish per client
    - qos: Quality of service level (0, 1, 2)
    - message: The message payload
    - topic: Topic to publish to
    - retain: Whether to retain the last message
    - clean: Whether to use a clean session
    - keepalive: Keepalive interval in seconds`,
	Run: func(cmd *cobra.Command, args []string) {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		// Parse flags
		clientID, _ := cmd.Flags().GetString("clientID")
		clients, _ := cmd.Flags().GetInt("clients")
		delay, _ := cmd.Flags().GetInt("delay")
		count, _ := cmd.Flags().GetInt("count")
		retain, _ := cmd.Flags().GetBool("retain")
		message, _ := cmd.Flags().GetString("message")
		topic, _ := cmd.Flags().GetString("topic")
		qos, _ := cmd.Flags().GetUint16("qos")
		cleanSession, _ := cmd.Flags().GetBool("clean")
		keepalive, _ := cmd.Flags().GetUint16("keepalive")

		b, err := bench.NewBenchmark(
			Cfg,
			bench.WithClientID(clientID),
			bench.WithClients(clients),
			bench.WithTopic(topic),
			bench.WithQoS(qos),
			bench.WithMessageCount(count),
			bench.WithDelay(delay),
			bench.WithRetained(retain),
			bench.WithCleanSession(cleanSession),
			bench.WithKeepAlive(keepalive),
			bench.WithMessage(message),
		)
		if err != nil {
			logger.Error("Failed to create benchmark", logger.ErrorAttr(err))
			return
		}

		go func() {
			<-sigs
			logger.Info("Received shutdown signal", logger.State("completed"))
			os.Exit(0)
		}()

		b.PublishMessages()
	},
}

func init() {
	rootCmd.AddCommand(pubCmd)

	pubCmd.Flags().StringP("clientID", "i", "benchmq-client", "Client ID for MQTT connections")
	pubCmd.Flags().IntP("clients", "c", 100, "Number of concurrent clients to connect")
	pubCmd.Flags().IntP("delay", "d", 1000, "Delay between messages in milliseconds")
	pubCmd.Flags().IntP("count", "n", 1000, "Number of messages to publish per client")
	pubCmd.Flags().BoolP("retain", "r", false, "Retain the last message")
	pubCmd.Flags().Uint16P("qos", "q", 0, "Quality of service level (0, 1, 2)")
	pubCmd.Flags().StringP("message", "m", "Hello, World!", "Message to publish")
	pubCmd.Flags().StringP("topic", "t", "benchmq", "Topic to publish messages to")
	pubCmd.Flags().BoolP("clean", "x", true, "Clean previous session when connecting")
	pubCmd.Flags().Uint16P("keepalive", "k", 60, "Keepalive interval in seconds")
}
