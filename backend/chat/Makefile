PORT_COMMAND ?= 8080
PORT_QUERY ?= 8081
PORT_WS ?= 8082
PF_DIR ?= .tmp/port-forward

.PHONY: run-command run-query run-ws run-consumer-notifier run-consumer-chat run-backend \
	tidy-go-mod go-deps-update \
	docker-build-api-command docker-build-api-query docker-build-ws docker-build-consumer-notifier docker-build-consumer-chat docker-build-all \
	load-orb-api-command load-orb-api-query load-orb-ws load-orb-consumer-notifier load-orb-consumer-chat load-orb-all docker-build-and-load-all \
	k8s-apply k8s-delete k8s-restart k8s-status k8s-logs k8s-deploy k8s-endpoints \
	k8s-port-forward-start k8s-port-forward-stop k8s-port-forward-status

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

.PHONY: ensure-orb-image load-orb-api-command load-orb-api-query load-orb-ws load-orb-consumer-notifier load-orb-consumer-chat load-orb-all docker-build-and-load-all \
	k8s-apply k8s-delete k8s-restart k8s-status k8s-logs k8s-deploy k8s-endpoints \
	k8s-port-forward-start k8s-port-forward-stop k8s-port-forward-status

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

k8s-apply:
	kubectl apply -f k8s/deployment.yaml
	kubectl apply -f k8s/ingress.yaml

k8s-delete:
	kubectl delete -f k8s/ingress.yaml --ignore-not-found
	kubectl delete -f k8s/deployment.yaml --ignore-not-found

k8s-restart:
	kubectl rollout restart deploy/chat-api-command deploy/chat-api-query deploy/chat-ws deploy/chat-consumer-notifier deploy/chat-consumer-chat

k8s-status:
	kubectl get deploy
	kubectl get pods
	kubectl get svc
	kubectl get ingress

k8s-logs:
	kubectl logs deploy/chat-api-command --tail=100
	kubectl logs deploy/chat-api-query --tail=100
	kubectl logs deploy/chat-ws --tail=100
	kubectl logs deploy/chat-consumer-notifier --tail=100
	kubectl logs deploy/chat-consumer-chat --tail=100

k8s-endpoints:
	@echo "query API (GET):    http://chat.local/api/chats/..."
	@echo "command API (POST): http://chat-command.local/api/chats/..."
	@echo "websocket:          ws://chat.local/ws"
	kubectl get ingress

k8s-port-forward-start:
	mkdir -p $(PF_DIR)
	@if [ -f $(PF_DIR)/api-command.pid ] && kill -0 $$(cat $(PF_DIR)/api-command.pid) 2>/dev/null; then \
		echo "api-command port-forward already running"; \
	else \
		kubectl port-forward svc/chat-api-command 8080:8080 > $(PF_DIR)/api-command.log 2>&1 & echo $$! > $(PF_DIR)/api-command.pid; \
		echo "started api-command port-forward (8080)"; \
	fi
	@if [ -f $(PF_DIR)/api-query.pid ] && kill -0 $$(cat $(PF_DIR)/api-query.pid) 2>/dev/null; then \
		echo "api-query port-forward already running"; \
	else \
		kubectl port-forward svc/chat-api-query 8081:8081 > $(PF_DIR)/api-query.log 2>&1 & echo $$! > $(PF_DIR)/api-query.pid; \
		echo "started api-query port-forward (8081)"; \
	fi
	@if [ -f $(PF_DIR)/ws.pid ] && kill -0 $$(cat $(PF_DIR)/ws.pid) 2>/dev/null; then \
		echo "ws port-forward already running"; \
	else \
		kubectl port-forward svc/chat-ws 8082:8082 > $(PF_DIR)/ws.log 2>&1 & echo $$! > $(PF_DIR)/ws.pid; \
		echo "started ws port-forward (8082)"; \
	fi

k8s-port-forward-stop:
	@if [ -f $(PF_DIR)/api-command.pid ] && kill -0 $$(cat $(PF_DIR)/api-command.pid) 2>/dev/null; then \
		kill $$(cat $(PF_DIR)/api-command.pid) && echo "stopped api-command port-forward"; \
	else \
		echo "api-command port-forward not running"; \
	fi
	@if [ -f $(PF_DIR)/api-query.pid ] && kill -0 $$(cat $(PF_DIR)/api-query.pid) 2>/dev/null; then \
		kill $$(cat $(PF_DIR)/api-query.pid) && echo "stopped api-query port-forward"; \
	else \
		echo "api-query port-forward not running"; \
	fi
	@if [ -f $(PF_DIR)/ws.pid ] && kill -0 $$(cat $(PF_DIR)/ws.pid) 2>/dev/null; then \
		kill $$(cat $(PF_DIR)/ws.pid) && echo "stopped ws port-forward"; \
	else \
		echo "ws port-forward not running"; \
	fi

k8s-port-forward-status:
	@if [ -f $(PF_DIR)/api-command.pid ] && kill -0 $$(cat $(PF_DIR)/api-command.pid) 2>/dev/null; then \
		echo "api-command: running (http://localhost:8080)"; \
	else \
		echo "api-command: stopped"; \
	fi
	@if [ -f $(PF_DIR)/api-query.pid ] && kill -0 $$(cat $(PF_DIR)/api-query.pid) 2>/dev/null; then \
		echo "api-query: running (http://localhost:8081)"; \
	else \
		echo "api-query: stopped"; \
	fi
	@if [ -f $(PF_DIR)/ws.pid ] && kill -0 $$(cat $(PF_DIR)/ws.pid) 2>/dev/null; then \
		echo "ws: running (ws://localhost:8082/ws)"; \
	else \
		echo "ws: stopped"; \
	fi

# build → image確認 → manifest適用 → rollout restart → status まで一括
k8s-deploy: docker-build-and-load-all k8s-apply k8s-restart k8s-status

.PHONY: compose-up compose-up-build compose-down compose-logs

compose-up:
	cd .. && docker compose -f compose.yml up -d

compose-up-build:
	cd .. && docker compose -f compose.yml up -d --build

compose-down:
	cd .. && docker compose -f compose.yml down

compose-logs:
	cd .. && docker compose -f compose.yml logs -f --tail=100

