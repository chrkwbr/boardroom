z_run-command:
	cd backend/chat && PORT=8080 go run ./cmd/api-command/main.go

z_run-query:
	cd backend/chat && PORT=8081 go run ./cmd/api-query/main.go

z_run-ws:
	cd backend/chat && PORT=8082 go run ./cmd/ws/main.go

z_run-consumer-notifier:
	cd backend/chat && go run ./cmd/consumer-notifier/main.go

z_run-consumer-chat:
	cd backend/chat && go run ./cmd/consumer-chat/main.go


run-backend:
	trap 'kill 0' SIGINT SIGTERM EXIT; \
	($(MAKE) z_run-command 2>&1 | sed 's/^/[command-api] /') & \
	($(MAKE) z_run-query 2>&1 | sed 's/^/[query-api]  /') & \
	($(MAKE) z_run-ws 2>&1 | sed 's/^/[ws]         /') & \
	($(MAKE) z_run-consumer-notifier 2>&1 | sed 's/^/[notifier]  /') & \
	($(MAKE) z_run-consumer-chat 2>&1 | sed 's/^/[consumer]  /') & \
	wait

migrate-scylla:
	cd appinfra && docker compose up -d scylla
	cd appinfra && docker compose run --rm --no-deps scylla-init

# ローカルに cqlsh を入れている場合のみ使用
migrate-scylla-local:
	MIGRATIONS_DIR=appinfra/scylla/migrations appinfra/scylla/migrate.sh

# chat keyspace を削除して、マイグレーションを 0 から再適用
reset-scylla:
	cd appinfra && docker compose up -d scylla
	cd appinfra && docker compose exec -T scylla cqlsh -e "DROP KEYSPACE IF EXISTS chat;"
	$(MAKE) migrate-scylla

kill-backend:
	pkill -f "cmd/chat" || true
	pkill -f "go-build.*/exe/main" || true