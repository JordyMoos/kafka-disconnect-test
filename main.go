package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Waiting 10 sec for kafka to come up in docker")
	time.Sleep(10 * time.Second)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:         []string{"kafka:9092"},
		Topic:           "example",
		Balancer:        &kafka.LeastBytes{},
		MaxAttempts:     10,
		IdleConnTimeout: 10 * time.Second, // 9 * time.Minutes
		Async:           false,
	})
	defer func() {
		fmt.Println("Closing kafka")
		writer.Close()
		fmt.Println("Kafka closed")

		fmt.Println("Finished!")
	}()

	var (
		err error
		ctx context.Context
	)

	// Context to stop after 10 seconds
	// Write the two messages to kafka
	// This block should work
	batchA := []kafka.Message{
		kafka.Message{Value: []byte("message-1")},
		kafka.Message{Value: []byte("message-2")},
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = writer.WriteMessages(ctx, batchA...)
	if err != nil {
		fmt.Println("Crashing here should not happen for the demo!!!")
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("")
	fmt.Println("Now: stop kafka")
	fmt.Println("Type: 'docker-compose stop kafka' in the terminal")
	fmt.Println("")
	fmt.Println("Sleeping for 30 sec...")
	time.Sleep(30 * time.Second)
	fmt.Println("Continues")

	// Now we write again which will crash
	batchB := []kafka.Message{
		kafka.Message{Value: []byte("message-1")},
		kafka.Message{Value: []byte("message-2")},
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = writer.WriteMessages(ctx, batchB...)
	if err != nil {
		fmt.Println("The write is crashed, which is expected")
		fmt.Println("Error: ", err)
		fmt.Println("Kafka should Close now")
		return
	}
}
