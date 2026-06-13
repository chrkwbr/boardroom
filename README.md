# Boardroom Kubernetes Runbook

OrbStack + kubectl で chat アプリを起動するための最小手順です。

## 1. 前提確認

```bash
# OrbStack の Kubernetes クラスタを起動
orbctl start k8s

# いま kubectl が向いているクラスタ(context)を確認
kubectl config current-context
```

`current-context` が `orbstack` になっていることを確認します。

## 2. イメージをビルド

```bash
cd /Users/chiharu-kuwabara/dev/boardroom

# chat の各サービスイメージをローカル Docker に build
make build-chat-images
```

OrbStack は Docker エンジン一体型のため、ローカルで build したイメージをそのまま Kubernetes で使えます。

## 3. イメージ存在チェック（任意）

```bash
cd /Users/chiharu-kuwabara/dev/boardroom

# k8s で使う予定のローカルイメージが存在するか確認
make orb-load-all
```

このターゲットは `boardroom/*:latest` イメージの存在確認を行います。

## 4. Deployment 適用

```bash
cd /Users/chiharu-kuwabara/dev/boardroom

# deployment/service/configmap をクラスタへ反映（作成または更新）
kubectl apply -f backend/chat/k8s/deployment.yaml
```

## 5. 状態確認

```bash
# Deployment の希望状態/現在状態を確認
kubectl get deploy

# Pod の起動状態をリアルタイム監視
kubectl get pods -w

# Service の公開ポート・到達先を確認
kubectl get svc
```

## 6. ログ確認

```bash
# 各 Deployment 配下 Pod の標準出力ログを追跡
kubectl logs deploy/chat-api-command -f
kubectl logs deploy/chat-api-query -f
kubectl logs deploy/chat-ws -f
kubectl logs deploy/chat-consumer-notifier -f
kubectl logs deploy/chat-consumer-chat -f
```

## 7. ローカルから API を叩く（ポートフォワード）

```bash
# クラスタ内 Service のポートをローカルへトンネル
kubectl port-forward svc/chat-api-command 8080:8080
kubectl port-forward svc/chat-api-query 8081:8081
kubectl port-forward svc/chat-ws 8082:8082
```

## 8. 再起動・削除

```bash
# 各 Deployment をローリング再起動（設定変更反映時など）
kubectl rollout restart deploy/chat-api-command
kubectl rollout restart deploy/chat-api-query
kubectl rollout restart deploy/chat-ws
kubectl rollout restart deploy/chat-consumer-notifier
kubectl rollout restart deploy/chat-consumer-chat
```

```bash
# マニフェストで作成したリソース一式を削除
kubectl delete -f backend/chat/k8s/deployment.yaml
```

## 9. 接続先を変更したい場合

`backend/chat/k8s/deployment.yaml` の `ConfigMap` (`chat-app-config`) で以下を変更します。

- `KAFKA_BROKERS`
- `REDIS_ADDR`
- `SCYLLA_HOST`

変更後は再適用します。

```bash
# 変更済みマニフェストを再反映
kubectl apply -f backend/chat/k8s/deployment.yaml
```
