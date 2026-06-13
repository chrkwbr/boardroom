module chat-consumer-chat

go 1.25.4

require (
	boardroom/chat-readmodel v0.0.0
	boardroom/shared v0.0.0
)

require (
	github.com/gocql/gocql v1.7.0 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/pierrec/lz4/v4 v4.1.27 // indirect
	github.com/segmentio/kafka-go v0.4.51 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

replace (
	boardroom/chat-readmodel => ../../internal/readmodel
	boardroom/shared => ../../../pkg/shared
)
