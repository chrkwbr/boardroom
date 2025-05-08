package processor

import (
	"backend/infra/hub"
	"backend/infra/pubsub/kafka"
	"log"
)

type KafkaWsProcessor struct {
}

func NewKafkaWsProcessor() *KafkaWsProcessor {
	return &KafkaWsProcessor{}
}

func (k *KafkaWsProcessor) Start() {
	chat_event_kafka, err := hub.GetHubFactory().GetHub(hub.ChatEventKafkaWs)
	if err != nil {
		log.Panicln("Failed to get hub:", err)
	}
	sub := kafka.NewKafkaReader([]string{"localhost:9092"}, "chat_messages")
	go func() {
		if err := sub.Subscribe("_", func(key string, value []byte) error {
			chat_event_kafka.BroadcastMessage(value)
			return nil
		}); err != nil {
			panic(err)
		}
	}()

}
