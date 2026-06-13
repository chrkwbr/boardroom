PORT_COMMAND ?= 8080
PORT_QUERY ?= 8081
PORT_WS ?= 8082

.PHONY: run-command run-query run-ws run-consumer-notifier run-consumer-chat run-backend \
	tidy-go-mod go-deps-update \
	docker-build-api-command docker-build-api-query docker-build-ws docker-build-consumer-notifier docker-build-consumer-chat docker-build-all \
	load-orb-api-command load-orb-api-query load-orb-ws load-orb-consumer-notifier load-orb-consumer-chat load-orb-all docker-build-and-load-all

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
	docker build -f cmd/api-command/Dockerfile -t boardroom/api-command:latest ../..

docker-build-api-query:
	docker build -f cmd/api-query/Dockerfile -t boardroom/api-query:latest ../..

docker-build-ws:
	docker build -f cmd/ws/Dockerfile -t boardroom/ws:latest ../..

docker-build-consumer-notifier:
	docker build -f cmd/consumer-notifier/Dockerfile -t boardroom/consumer-notifier:latest ../..

docker-build-consumer-chat:
	docker build -f cmd/consumer-chat/Dockerfile -t boardroom/consumer-chat:latest ../..

docker-build-all: \
	docker-build-api-command \
	docker-build-api-query \
	docker-build-ws \
	docker-build-consumer-notifier \
	docker-build-consumer-chat

.PHONY: ensure-orb-image load-orb-api-command load-orb-api-query load-orb-ws load-orb-consumer-notifier load-orb-consumer-chat load-orb-all docker-build-and-load-all

ensure-orb-image:
	@docker image inspect $(IMAGE) >/dev/null 2>&1 || (echo "[error] image not found: $(IMAGE)"; echo "        run: make docker-build-all"; exit 1)

load-orb-api-command: IMAGE=boardroom/api-command:latest
load-orb-api-command: ensure-orb-image
	@echo "[orb] using local image: $(IMAGE)"

load-orb-api-query: IMAGE=boardroom/api-query:latest
load-orb-api-query: ensure-orb-image
	@echo "[orb] using local image: $(IMAGE)"

load-orb-ws: IMAGE=boardroom/ws:latest
load-orb-ws: ensure-orb-image
	@echo "[orb] using local image: $(IMAGE)"

load-orb-consumer-notifier: IMAGE=boardroom/consumer-notifier:latest
load-orb-consumer-notifier: ensure-orb-image
	@echo "[orb] using local image: $(IMAGE)"

load-orb-consumer-chat: IMAGE=boardroom/consumer-chat:latest
load-orb-consumer-chat: ensure-orb-image
	@echo "[orb] using local image: $(IMAGE)"

load-orb-all: \
	load-orb-api-command \
	load-orb-api-query \
	load-orb-ws \
	load-orb-consumer-notifier \
	load-orb-consumer-chat
	@echo "[orb] image check complete. Deploy with imagePullPolicy: IfNotPresent or Never."

docker-build-and-load-all: docker-build-all load-orb-all
