package event

type EventPublisher interface {
	Publish(topic string, key string, value []byte) error
}

type EventSubscriber interface {
	Subscribe(topic string, handler func(key string, value []byte) error) error
}
