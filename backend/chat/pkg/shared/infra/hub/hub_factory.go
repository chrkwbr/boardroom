package hub

import (
	"fmt"
	"sync"
)

type hubKey string

const (
	ChatEventOutbox   hubKey = "chat_event_outbox"
	ChatEventWsPusher hubKey = "chat_event_ws_pusher"
)

var hubParams = map[hubKey]*HubParams{
	ChatEventOutbox: {
		broadcastBuffer:  256,
		registerBuffer:   10,
		unregisterBuffer: 10,
	},
	ChatEventWsPusher: {
		broadcastBuffer:  256,
		registerBuffer:   256,
		unregisterBuffer: 256,
	},
}

type HubParams struct {
	broadcastBuffer  int32
	registerBuffer   int32
	unregisterBuffer int32
}

type HubFactory struct {
	mu   sync.Mutex
	hubs map[hubKey]*Hub
}

var (
	once     sync.Once
	instance *HubFactory
)

func GetHubFactory() *HubFactory {
	once.Do(func() {
		instance = &HubFactory{
			hubs: make(map[hubKey]*Hub),
		}
	})
	return instance
}

func (f *HubFactory) SetNewHubFn(fn func() *Hub) {
	f.mu.Lock()
	defer f.mu.Unlock()
}

func (f *HubFactory) GetHub(name hubKey) (*Hub, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if hub, exists := f.hubs[name]; exists {
		return hub, nil
	}

	hubParam := hubParams[name]
	if hubParam == nil {
		return nil, fmt.Errorf("hub params not found for %s", name)
	}
	hub := newHub(hubParam.broadcastBuffer, hubParam.registerBuffer, hubParam.unregisterBuffer)
	f.hubs[name] = hub

	go hub.Run()
	return hub, nil
}
