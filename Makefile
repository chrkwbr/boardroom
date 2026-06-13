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

# chat image をOrbStack k8sへload
load-chat-images:
	$(CHAT_MAKE) load-orb-all

build-and-load-chat-images:
	$(CHAT_MAKE) docker-build-and-load-all

# 互換エイリアス
orb-load-all: load-chat-images

k8s-apply:
	$(CHAT_MAKE) k8s-apply

k8s-delete:
	$(CHAT_MAKE) k8s-delete

k8s-restart:
	$(CHAT_MAKE) k8s-restart

k8s-status:
	$(CHAT_MAKE) k8s-status

k8s-logs:
	$(CHAT_MAKE) k8s-logs

k8s-endpoints:
	$(CHAT_MAKE) k8s-endpoints

k8s-port-forward-start:
	$(CHAT_MAKE) k8s-port-forward-start

k8s-port-forward-stop:
	$(CHAT_MAKE) k8s-port-forward-stop

k8s-port-forward-status:
	$(CHAT_MAKE) k8s-port-forward-status

# build + apply + restart + status
k8s-deploy:
	$(CHAT_MAKE) k8s-deploy

compose-backend-up:
	$(CHAT_MAKE) compose-up

compose-backend-up-build:
	$(CHAT_MAKE) compose-build-up

compose-backend-down:
	$(CHAT_MAKE) compose-down

compose-backend-logs:
	$(CHAT_MAKE) compose-logs
