z_run-command:
	cd backend && PORT=8080 go run ./cmd/chat/api/command/main.go

z_run-query:
	cd backend && PORT=8081 go run ./cmd/chat/api/query/main.go

z_run-ws:
	cd backend && PORT=8082 go run ./cmd/chat/ws/main.go

z_run-consumer-notifier:
	cd backend && go run ./cmd/chat/consumer/notifier/main.go

z_run-consumer-chat:
	cd backend && go run ./cmd/chat/consumer/chat/main.go


run-backend:
	trap 'kill 0' SIGINT SIGTERM EXIT; \
	($(MAKE) z_run-command 2>&1 | sed 's/^/[command-api] /') & \
	($(MAKE) z_run-query 2>&1 | sed 's/^/[query-api]  /') & \
	($(MAKE) z_run-ws 2>&1 | sed 's/^/[ws]         /') & \
	($(MAKE) z_run-consumer-notifier 2>&1 | sed 's/^/[notifier]  /') & \
	($(MAKE) z_run-consumer-chat 2>&1 | sed 's/^/[consumer]  /') & \
	wait

kill-backend:
	pkill -f "cmd/chat" || true
	pkill -f "go-build.*/exe/main" || true