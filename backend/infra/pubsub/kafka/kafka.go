package kafka

import (
	"backend/infra/pubsub"
	"context"
	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

type KafkaReader struct {
	reader *kafka.Reader
}

func NewKafkaWriter(brokers []string) pubsub.EventPublisher {
	return &KafkaWriter{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
			Async:    true,
		},
	}
}

func NewKafkaReader(brokers []string, topic string, groupId string) pubsub.EventSubscriber {
	return &KafkaReader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupId,
			MinBytes: 1,    // ToDo: adjust as needed
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (kw *KafkaWriter) Publish(topic string, key string, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}
	return kw.writer.WriteMessages(context.Background(), msg)
}

func (k KafkaReader) Subscribe(topic string, handler func(key string, value []byte) error) error {
	for {
		msg, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}
		if err := handler(string(msg.Key), msg.Value); err != nil {
			return err
		}
	}
}
