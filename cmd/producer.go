/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

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

		if _, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: []byte("hello"),
		}); err != nil {
			log.Println("Failed to send message", err)
			return
		}
		defer producer.Close()

		if err != nil {
			log.Println("Failed to publish message", err)
		} else {
			log.Println("Published message")
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
