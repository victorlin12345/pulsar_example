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
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/spf13/cobra"
)

// producerCmd represents the producer command
var producerCmd = &cobra.Command{
	Use:   "producer",
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
			log.Println("Failed to create Pulsar client", err)
			return
		}
		defer client.Close()

		producer, err := client.CreateProducer(pulsar.ProducerOptions{
			Topic: "investments/stocks/stock-ticker",
		})
		if err != nil {
			log.Println("Failed to create producer", err)
			return
		}
		defer producer.Close()

		// Create a ticker that ticks every 1 second
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		// Create a context with cancellation
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Set up signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Counter for messages
		messageCount := 0
	ProcessLoop:
		for {
			select {
			case <-ticker.C:
				messageCount++
				key := fmt.Sprintf("%d", messageCount%3)
				payload := fmt.Sprintf("hello %d", messageCount%3)
				if _, err = producer.Send(ctx, &pulsar.ProducerMessage{
					Key:     key,
					Payload: []byte(payload),
				}); err != nil {
					log.Println("Failed to send message", err)
					return
				}
				log.Printf("produce message key:%s value:%s\n", key, payload)

			case <-sigChan:
				log.Println("Shutting down producer...")
				break ProcessLoop
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(producerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// producerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// producerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
