# ローカル開発環境 構築手順 (Kubernetes / Helm 編)

手元の Mac (OrbStack) 上の Kubernetes クラスター内に、Go アプリケーションが使用するミドルウェア（Redis, Kafka, ScyllaDB 互換コンテナ）を Helm を用いて一発で構築する手順です。

すべてのミドルウェアは Kubernetes 内部のプライベート DNS で名前解決され、Go アプリの各プロセスから透過的に接続可能になります。

## 前提条件

事前に手元の Mac に以下のツールがインストールされ、有効化されていることを確認してください。

- OrbStack: Kubernetes 機能を有効化（Preferences -> Kubernetes -> Enable Kubernetes）
- Homebrew: パッケージマネージャー

[コマンド]
```
brew install helm kubectl
```

## ミドルウェアの構築手順

Kubernetes 内に開発用の各コンテナをスタンドアロン（単一ノード）、認証なしの最小構成で立ち上げます。

### Step 1: Helm レポジトリの追加 & 更新
共通で使用する Bitnami の公式チャートレポジトリを追加します。

[コマンド]
```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

### Step 2: Redis (Pub/Sub・WebSocket用) のインストール

リアルタイム通知のブロードキャストに使用する Redis を、メモリ消費の少ないスタンドアロンモードで起動します。

[コマンド]
```
helm install boardroom-redis bitnami/redis \
  --set architecture=standalone \
  --set auth.enabled=false
```

### Step 3: Kafka (イベント駆動用) のインストール
メッセージのイベント駆動に使用する Kafka を、Zookeeper 不要の KRaft モード・単一ブローカーで起動します。

[コマンド]
```
helm install boardroom-kafka bitnami/kafka \
  --set listeners.client.protocol=PLAINTEXT \
  --set sasl.enabled=false \
  --set rbac.create=false \
  --set broker.replicaCount=1 \
  --set persistence.enabled=false \
  --set volumePermissions.enabled=true
```

### Step 4: ScyllaDB / Cassandra (永続ストレージ用) のインストール
ローカル開発時のマシンリソースを節約するため、ScyllaDB と完全なプロトコル互換性を持つ Cassandra チャートを 1 ノードで流用します。Go 側の CQL ドライバ（gocql等）はそのまま 100% 動作します。

[コマンド]
```
helm install boardroom-scylla bitnami/cassandra \
  --set replicaCount=1 \
  --set dbUser.enabled=false
  --set image.registry=quay.io \
  --set image.repository=bitnami/cassandra \
  --set image.tag=5.0.0
```

## 起動ステータスの確認

すべての Pod が正常に起動し、STATUS が Running になるまで数分待ちます。

[コマンド]
```
kubectl get pods -w
```

[正常に起動した際の見本]
```
NAME                 READY   STATUS    RESTARTS   AGE
boardroom-kafka-0           1/1     Running   0          2m
boardroom-redis-master-0    1/1     Running   0          3m
boardroom-scylla-0          1/1     Running   0          1m
```

## Go アプリケーションからの接続設定 (接続先一覧)

Kubernetes 内にデプロイする Go アプリの各プロセス（Deployment マニフェストなど）の環境変数には、以下の内部ドメイン（DNS）を指定して接続してください。

- Redis
接続先: boardroom-redis-master.default.svc.cluster.local:6379
担当プロセス: chat-ws, consumer-notifier

- Kafka
接続先: boardroom-kafka.default.svc.cluster.local:9092
担当プロセス: chat-command, consumer-chat, consumer-notifier

- ScyllaDB
接続先: boardroom-scylla.default.svc.cluster.local:9042
担当プロセス: chat-query, consumer-chat

## 環境の削除・初期化

ローカルの検証環境を完全に削除して一からやり直したい場合は、以下のコマンドを実行してください。永続ボリューム（データ実体）も含めてすべてクリーンアップされます。

[コマンド]
```
helm uninstall boardroom-redis my-kafka my-scylla
kubectl delete pvc --all
```