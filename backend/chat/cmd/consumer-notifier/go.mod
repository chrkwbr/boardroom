module chat-consumer-notifier

go 1.25.4

require (
	boardroom/chat-shared v0.0.0
	github.com/redis/go-redis/v9 v9.20.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/gocql/gocql v1.7.0 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/segmentio/kafka-go v0.4.51 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

replace boardroom/chat-shared => ../../../pkg/shared
