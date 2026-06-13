package main

import (
	"boardroom/chat-readmodel"
	"boardroom/shared/infra/pubsub"
	"chat-consumer-chat/internal"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func kafkaBrokers() []string {
	brokers := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	if brokers == "" {
		return []string{"localhost:9092"}
	}
	parts := strings.Split(brokers, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return []string{"localhost:9092"}
	}
	return out
}

func scyllaHosts() []string {
	hosts := strings.TrimSpace(os.Getenv("SCYLLA_HOST"))
	if hosts == "" {
		return []string{"localhost"}
	}
	parts := strings.Split(hosts, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return []string{"localhost"}
	}
	return out
}

func main() {
	log.Println("==== Starting consumer-kafka-chat...")

	scylla, err := readmodel.NewChatScyllaRepository(scyllaHosts()...)
	if err != nil {
		log.Fatalf("Failed to connect to ScyllaDB: %v", err)
	}
	defer scylla.Close()

	kafkaReader := pubsub.NewKafkaReader(kafkaBrokers(), "chat-events", "chat-materializer")
	defer kafkaReader.Close()

	internal.NewMaterializer(kafkaReader, scylla).Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("===== Shutting down consumer-kafka-chat...")
}
