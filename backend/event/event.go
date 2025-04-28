package event

import "sync"

type Event struct {
	Name string
	Data interface{}
}

var (
	subscribers = make(map[string][]chan Event)
	mu          sync.RWMutex
)

func Subscribe(eventName string) chan Event {
	mu.Lock()
	defer mu.Unlock()

	ch := make(chan Event)
	subscribers[eventName] = append(subscribers[eventName], ch)
	return ch
}

func Publish(eventName string, data interface{}) {
	mu.RLock()
	subs := make([]chan Event, 0)
	if channels, ok := subscribers[eventName]; ok {
		subs = append(subs, channels...)
	}
	mu.RUnlock()

	event := Event{Name: eventName, Data: data}
	for _, ch := range subs {
		go func(c chan Event) {
			defer func() {
				// 閉じられたチャネルへの送信によるパニックを回避
				if r := recover(); r != nil {
					// ここでログ出力などが可能
				}
			}()
			c <- event
		}(ch)
	}
}

func Unsubscribe(eventName string, ch chan Event) {
	mu.Lock()
	defer mu.Unlock()

	if subs, ok := subscribers[eventName]; ok {
		for i, subscriber := range subs {
			if subscriber == ch {
				// チャネルを購読リストから削除
				subscribers[eventName] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}
}
