#!/bin/bash
set -e

HOST="${SCYLLA_HOST:-localhost}"
PORT="${SCYLLA_PORT:-9042}"
KEYSPACE="chat"
MIGRATIONS_DIR="${MIGRATIONS_DIR:-$(dirname "$0")/migrations}"

# keyspace と migrations テーブルの作成
cqlsh "$HOST" "$PORT" <<EOF
CREATE KEYSPACE IF NOT EXISTS $KEYSPACE
    WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE IF NOT EXISTS $KEYSPACE.schema_migrations (
    version text PRIMARY KEY,
    applied_at timestamp
);
EOF

echo "Checking migrations in $MIGRATIONS_DIR ..."

for file in $(ls "$MIGRATIONS_DIR"/*.cql 2>/dev/null | sort); do
    version=$(basename "$file" .cql)

    # 適用済みかチェック
    applied=$(cqlsh "$HOST" "$PORT" -e \
        "SELECT version FROM $KEYSPACE.schema_migrations WHERE version = '$version';" \
        | grep -c "$version" || true)

    if [ "$applied" -gt 0 ]; then
        echo "  [skip] $version (already applied)"
        continue
    fi

    echo "  [apply] $version ..."
    cqlsh "$HOST" "$PORT" -k "$KEYSPACE" -f "$file"

    cqlsh "$HOST" "$PORT" -e \
        "INSERT INTO $KEYSPACE.schema_migrations (version, applied_at) VALUES ('$version', toTimestamp(now()));"

    echo "  [done]  $version"
done

echo "Migration complete."


