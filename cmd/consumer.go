/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/spf13/cobra"
)

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pulsar.NewClient(pulsar.ClientOptions{
			URL: "pulsar://localhost:6650",
		})
		if err != nil {
			log.Printf("Failed to create Pulsar client\n:%s", err)
			return
		}

		defer client.Close()

		// Create a context with cancellation
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Set up signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		consumer, err := client.Subscribe(pulsar.ConsumerOptions{
			Topic:            "investments/stocks/stock-ticker",
			SubscriptionName: "my-sub",
			Type:             pulsar.KeyShared,
		})
		if err != nil {
			log.Printf("Failed to create consumer\n:%s", err)
			return
		}
		defer consumer.Close()

		// Start continuous message consumption
	ProcessLoop:
		for {
			select {
			case <-sigChan:
				log.Println("Shutting down consumer...")
				break ProcessLoop
			default:
				msg, err := consumer.Receive(ctx)
				if err != nil {
					log.Printf("Error receiving message: %v", err)
					continue
				}

				// Process the message
				fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
					msg.ID(), string(msg.Payload()))

				// Acknowledge the message
				consumer.Ack(msg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consumerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
