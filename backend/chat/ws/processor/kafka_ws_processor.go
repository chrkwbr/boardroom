package processor

import (
	"backend/infra/hub"
	"backend/infra/pubsub"
	"backend/infra/pubsub/kafka"
	"log"
)

type KafkaWsProcessor struct {
	subscriber pubsub.EventSubscriber
}

func NewKafkaWsProcessor() *KafkaWsProcessor {
	return &KafkaWsProcessor{
		subscriber: kafka.NewKafkaReader([]string{"localhost:9092"}, "chat_messages", "websocket_processor"),
	}
}

func (k *KafkaWsProcessor) Start() {
	chat_event_kafka, err := hub.GetHubFactory().GetHub(hub.ChatEventKafkaWs)
	if err != nil {
		log.Panicln("Failed to get hub:", err)
	}
	go func() {
		if err := k.subscriber.Subscribe("_", func(key string, value []byte) error {
			chat_event_kafka.BroadcastMessage(value)
			return nil
		}); err != nil {
			panic(err)
		}
	}()

}
