PORT_COMMAND ?= 8080
PORT_QUERY ?= 8081
PORT_WS ?= 8082

.PHONY: run-command run-query run-ws run-consumer-notifier run-consumer-chat run-backend \
	tidy-go-mod go-deps-update \
	docker-build-api-command docker-build-api-query docker-build-ws docker-build-consumer-notifier docker-build-consumer-chat docker-build-all

run-command:
	cd cmd/api-command && PORT=$(PORT_COMMAND) go run .

run-query:
	cd cmd/api-query && PORT=$(PORT_QUERY) go run .

run-ws:
	cd cmd/ws && PORT=$(PORT_WS) go run .

run-consumer-notifier:
	cd cmd/consumer-notifier && go run .

run-consumer-chat:
	cd cmd/consumer-chat && go run .

run-backend:
	trap 'kill 0' SIGINT SIGTERM EXIT; \
	($(MAKE) run-command 2>&1 | sed 's/^/[command-api] /') & \
	($(MAKE) run-query 2>&1 | sed 's/^/[query-api]  /') & \
	($(MAKE) run-ws 2>&1 | sed 's/^/[ws]         /') & \
	($(MAKE) run-consumer-notifier 2>&1 | sed 's/^/[notifier]  /') & \
	($(MAKE) run-consumer-chat 2>&1 | sed 's/^/[consumer]  /') & \
	wait

tidy-go-mod:
	cd internal/domain && go mod tidy
	cd internal/readmodel && go mod tidy
	cd internal/notification && go mod tidy
	cd cmd/api-command && go mod tidy
	cd cmd/api-query && go mod tidy
	cd cmd/ws && go mod tidy
	cd cmd/consumer-notifier && go mod tidy
	cd cmd/consumer-chat && go mod tidy

go-deps-update:
	cd internal/domain && go get -u ./... && go mod tidy
	cd internal/readmodel && go get -u ./... && go mod tidy
	cd internal/notification && go get -u ./... && go mod tidy
	cd cmd/api-command && go get -u ./... && go mod tidy
	cd cmd/api-query && go get -u ./... && go mod tidy
	cd cmd/ws && go get -u ./... && go mod tidy
	cd cmd/consumer-notifier && go get -u ./... && go mod tidy
	cd cmd/consumer-chat && go get -u ./... && go mod tidy

docker-build-api-command:
	docker build -f backend/chat/cmd/api-command/Dockerfile -t boardroom/api-command:latest ../..

docker-build-api-query:
	docker build -f backend/chat/cmd/api-query/Dockerfile -t boardroom/api-query:latest ../..

docker-build-ws:
	docker build -f backend/chat/cmd/ws/Dockerfile -t boardroom/ws:latest ../..

docker-build-consumer-notifier:
	docker build -f backend/chat/cmd/consumer-notifier/Dockerfile -t boardroom/consumer-notifier:latest ../..

docker-build-consumer-chat:
	docker build -f backend/chat/cmd/consumer-chat/Dockerfile -t boardroom/consumer-chat:latest ../..

docker-build-all: \
	docker-build-api-command \
	docker-build-api-query \
	docker-build-ws \
	docker-build-consumer-notifier \
	docker-build-consumer-chat

