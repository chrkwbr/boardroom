z_run-command:
	cd backend && PORT=8080 go run ./cmd/chat/api/command/main.go

z_run-query:
	cd backend && PORT=8081 go run ./cmd/chat/api/query/main.go

z_run-ws:
	cd backend && PORT=8082 go run ./cmd/chat/ws/main.go

z_run-consumer-kafka-chat:
	cd backend && go run ./cmd/chat/consumer/chat/main.go

run-backend:
	trap 'kill 0' SIGINT SIGTERM EXIT; \
	$(MAKE) z_run-command & \
	$(MAKE) z_run-query & \
	$(MAKE) z_run-ws & \
	$(MAKE) z_run-consumer-kafka-chat & \
	wait

kill-backend:
	pkill -f "cmd/chat" || true
	pkill -f "go-build.*/exe/main" || true