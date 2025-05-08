package processor

import (
	"backend/infra/hub"
	"backend/infra/pubsub"
	"log"
)

type WsChatEventPusher struct {
	subscriber pubsub.EventSubscriber
}

func NewWsChatEventPusher(sub pubsub.EventSubscriber) *WsChatEventPusher {
	return &WsChatEventPusher{
		subscriber: sub,
	}
}

func (k *WsChatEventPusher) Start() {
	wsChatEventPusherHub, err := hub.GetHubFactory().GetHub(hub.ChatEventWsPusher)
	if err != nil {
		log.Panicln("Failed to get hub:", err)
	}
	go func() {
		if err := k.subscriber.Subscribe("_", func(key string, value []byte) error {
			wsChatEventPusherHub.BroadcastMessage(value)
			return nil
		}); err != nil {
			panic(err)
		}
	}()

}
