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
	- clientID: Base client ID prefix (each client appends "-<n>")
    - clients: Number of concurrent clients
    - delay: Delay between messages in milliseconds
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
		defer signal.Stop(sigs)

		// Parse flags
		clientID, err := cmd.Flags().GetString("clientID")
		if err != nil {
			logger.Error("Failed to parse client ID", logger.ErrorAttr(err))
			return
		}

		clients, err := cmd.Flags().GetInt("clients")
		if err != nil {
			logger.Error("Failed to parse number of clients", logger.ErrorAttr(err))
			return
		}

		delay, err := cmd.Flags().GetInt("delay")
		if err != nil {
			logger.Error("Failed to parse delay", logger.ErrorAttr(err))
			return
		}

		count, err := cmd.Flags().GetInt("count")
		if err != nil {
			logger.Error("Failed to parse message count", logger.ErrorAttr(err))
			return
		}

		retain, err := cmd.Flags().GetBool("retain")
		if err != nil {
			logger.Error("Failed to parse retain flag", logger.ErrorAttr(err))
			return
		}

		message, err := cmd.Flags().GetString("message")
		if err != nil {
			logger.Error("Failed to parse message", logger.ErrorAttr(err))
			return
		}

		topic, err := cmd.Flags().GetString("topic")
		if err != nil {
			logger.Error("Failed to parse topic", logger.ErrorAttr(err))
			return
		}

		qos, err := cmd.Flags().GetUint16("qos")
		if err != nil {
			logger.Error("Failed to parse QoS", logger.ErrorAttr(err))
			return
		}

		cleanSession, err := cmd.Flags().GetBool("clean")
		if err != nil {
			logger.Error("Failed to parse clean session flag", logger.ErrorAttr(err))
			return
		}

		keepalive, err := cmd.Flags().GetUint16("keepalive")
		if err != nil {
			logger.Error("Failed to parse keepalive", logger.ErrorAttr(err))
			return
		}

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
			logger.Error("Failed to create benchmark", logger.State("failed"), logger.ErrorAttr(err))
			return
		}

		go func() {
			<-sigs
			logger.Info("Received shutdown signal", logger.State("interrupted"))
			os.Exit(0)
		}()

		b.PublishMessages()
	},
}

func init() {
	rootCmd.AddCommand(pubCmd)

	// Register flags
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
