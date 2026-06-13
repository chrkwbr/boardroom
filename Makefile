CHAT_MAKE := $(MAKE) -C backend/chat
PKG_MAKE := $(MAKE) -C backend/pkg

# chat は root からは一括操作のみ
run-chat:
	$(CHAT_MAKE) run-backend

# 互換エイリアス
run-backend: run-chat

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

tidy-go-mod:
	$(PKG_MAKE) tidy-go-mod
	$(CHAT_MAKE) tidy-go-mod

go-deps-update:
	$(PKG_MAKE) go-deps-update
	$(CHAT_MAKE) go-deps-update

# chat の Docker image はまとめてビルド
build-chat-images:
	$(CHAT_MAKE) docker-build-all

# 互換エイリアス
docker-build-all: build-chat-images
