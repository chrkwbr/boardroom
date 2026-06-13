module chat-consumer-notifier

go 1.25.4

require (
	boardroom/chat-notification v0.0.0
	boardroom/shared v0.0.0
	github.com/redis/go-redis/v9 v9.20.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/segmentio/kafka-go v0.4.51 // indirect
	go.uber.org/atomic v1.11.0 // indirect
)

replace (
	boardroom/chat-notification => ../../internal/notification
	boardroom/shared => ../../../pkg/shared
)
