package pubsub

type EventPublisher interface {
	Publish(topic string, key string, value []byte) error
	Close() error
}

type EventSubscriber interface {
	Subscribe(topic string, handler func(key string, value []byte) error) error
	Close() error
}
