# Backend Migration (Atlas)

This backend uses Atlas for PostgreSQL schema migrations (versioned mode).

## Files

| ファイル | 役割 |
|---|---|
| `atlas.hcl` | Atlas の環境設定 |
| `db/schema.sql` | **スキーマの最終あるべき姿**（編集し続ける source of truth） |
| `migrations/` | DB に適用される SQL の履歴（Atlas が自動生成） |
| `migrations/atlas.sum` | マイグレーションの整合性チェックファイル（自動管理） |

> `db/schema.sql` と `migrations/` は**役割が異なる**。  
> `schema.sql` は常に最新のあるべき状態を表し、`migrations/` はその変遷を履歴として積み上げる。

## Prerequisites

- Atlas CLI (`brew install ariga/tap/atlas`)
- Docker（Atlas が差分計算用の一時 dev DB として使用）
- `localhost:5433` で PostgreSQL が起動済み（`appinfra/compose.yaml` の postgres サービス）

## 運用フロー

### スキーマを変更したいとき

```
1. db/schema.sql を編集（カラム追加、テーブル追加 など）
        ↓
2. atlas migrate diff <name> --env local
   → schema.sql と migrations/ の差分を Atlas が計算し、新しい SQL ファイルを自動生成
        ↓
3. atlas migrate apply --env local
   → 未適用のファイルだけ DB に適用（適用済みはスキップされる）
```

### 具体例：`chat_events` にカラムを追加する場合

```bash
# 1. db/schema.sql に列を追加
#    例: body TEXT NOT NULL を chat.chat_events に追加

# 2. 差分マイグレーションを生成
atlas migrate diff add_body_to_chat_events --env local
# → migrations/202606XXXXXXXX_add_body_to_chat_events.sql が生成される

# 3. DB に適用
atlas migrate apply --env local
```

## よく使うコマンド

```bash
# 適用状況の確認
atlas migrate status --env local

# 未適用マイグレーションを適用
atlas migrate apply --env local

# schema.sql の変更から新しいマイグレーションを生成
atlas migrate diff <name> --env local

# atlas.sum を手動で再生成（ファイルを手動編集した場合）
atlas migrate hash --env local
```

