package hub

import (
	"errors"
	"sync"
)

type hubKey string

const (
	ChatEventOutbox hubKey = "chat_event_outbox"
	ChatEventKafka  hubKey = "chat_event_kafka"
)

type HubFactory struct {
	mu       sync.Mutex
	hubs     map[hubKey]*Hub
	newHubFn func() *Hub
}

var (
	once     sync.Once
	instance *HubFactory
)

func GetHubFactory() *HubFactory {
	once.Do(func() {
		instance = &HubFactory{
			hubs: make(map[hubKey]*Hub),
			newHubFn: func() *Hub {
				return newHub()
			},
		}
	})
	return instance
}

func (f *HubFactory) SetNewHubFn(fn func() *Hub) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.newHubFn = fn
}

func (f *HubFactory) GetHub(name hubKey) (*Hub, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if hub, exists := f.hubs[name]; exists {
		return hub, nil
	}

	if f.newHubFn() == nil {
		return nil, errors.New("hub creation function is not set")
	}

	hub := f.newHubFn()
	f.hubs[name] = hub

	go hub.Run()
	return hub, nil
}
