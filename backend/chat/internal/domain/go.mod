module boardroom/chat-domain

go 1.25.4

require (
	boardroom/shared v0.0.0
	github.com/google/uuid v1.6.0
)

replace boardroom/shared => ../../../pkg/shared
