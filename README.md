# Boardroom 起動手順（Kubernetes / OrbStack）

ローカルで chat アプリを起動するときに、毎回やることを順番でまとめています。

## 起動チェックリスト

- [ ] OrbStack の k8s が起動している
- [ ] `kubectl` context が `orbstack`
- [ ] `appinfra`（Kafka/Redis/Scylla など）が起動している
- [ ] chat イメージを build 済み
- [ ] k8s マニフェストを apply 済み
- [ ] frontend 用 `port-forward` が起動している

## 1. 事前準備

```bash
# k8s 起動
orbctl start k8s

# context 確認
kubectl config current-context
```

`orbstack` になっていることを確認します。

## 2. インフラ起動（Kafka / Redis / Scylla / Postgres）

```bash
cd /Users/chiharu-kuwabara/dev/boardroom/appinfra
docker compose up -d
```

## 3. （必要時）Scylla マイグレーション

```bash
cd /Users/chiharu-kuwabara/dev/boardroom
make migrate-scylla
```

## 4. アプリのビルドとデプロイ

```bash
cd /Users/chiharu-kuwabara/dev/boardroom

# build + image確認 + apply + rollout restart + status
make k8s-deploy
```

## 5. 状態確認

```bash
cd /Users/chiharu-kuwabara/dev/boardroom
make k8s-status
```

## 6. frontend から接続する（port-forward）

`frontend/vite.config.ts` の proxy は `localhost:8080/8081/8082` を向く前提です。先に `port-forward` を起動します。

```bash
cd /Users/chiharu-kuwabara/dev/boardroom
make k8s-port-forward-start
make k8s-port-forward-status
```

その後 frontend を起動します。

```bash
cd /Users/chiharu-kuwabara/dev/boardroom/frontend
deno task dev
```

## 7. ログ確認・再起動

```bash
cd /Users/chiharu-kuwabara/dev/boardroom

# 各 deployment のログ末尾を確認
make k8s-logs

# deployment を再起動
make k8s-restart
```

## 8. 停止/クリーンアップ

```bash
cd /Users/chiharu-kuwabara/dev/boardroom

# frontend向け port-forward 停止
make k8s-port-forward-stop

# k8s リソース削除
make k8s-delete
```

必要ならインフラも停止します。

```bash
cd /Users/chiharu-kuwabara/dev/boardroom/appinfra
docker compose down
```

## トラブルシュート

```bash
# Podの状態
kubectl get pods -o wide

# 特定サービスのログ
kubectl logs deploy/chat-api-query --tail=200
kubectl logs deploy/chat-consumer-chat --tail=200

# ingress の確認
kubectl get ingress -o wide
```

`chat-consumer-chat` が `localhost:9092` へ接続しようとしているログが出る場合は、Kafka の `advertised.listeners` を見直して `appinfra` 側を再起動してください。
