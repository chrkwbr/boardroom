package main

import (
	"backend/chat/cmd/consumer-chat/internal"
	"backend/chat/pkg/shared/infra/pubsub/kafka"
	"backend/chat/pkg/shared/readmodel"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("==== Starting consumer-kafka-chat...")

	scylla, err := readmodel.NewChatScyllaRepository("localhost")
	if err != nil {
		log.Fatalf("Failed to connect to ScyllaDB: %v", err)
	}
	defer scylla.Close()

	kafkaReader := kafka.NewKafkaReader([]string{"localhost:9092"}, "chat-events", "chat-materializer")
	defer kafkaReader.Close()

	internal.NewMaterializer(kafkaReader, scylla).Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("===== Shutting down consumer-kafka-chat...")
}
