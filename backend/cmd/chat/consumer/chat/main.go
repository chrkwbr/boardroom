package main

import (
	"backend/internal/chat/readmodel"
	"backend/internal/shared/infra/pubsub/kafka"
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

	readmodel.NewMaterializer(kafkaReader, scylla).Start()

	log.Println("==== Consumer started. Waiting for events...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("===== Shutting down consumer-kafka-chat...")
}
