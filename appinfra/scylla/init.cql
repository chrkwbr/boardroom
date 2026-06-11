CREATE KEYSPACE IF NOT EXISTS chat
    WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE IF NOT EXISTS chat.chat_messages (
    room_id   uuid,
    created_at bigint,
    id        uuid,
    sender_id  uuid,
    sender_name text,
    sender_icon text,
    message    text,
    version    bigint,
    updated_at bigint,
    PRIMARY KEY ((room_id), created_at, id)
) WITH CLUSTERING ORDER BY (created_at DESC, id ASC);

-- id 単体 / room_id + id での検索を ALLOW FILTERING なしに行うためのルックアップテーブル
CREATE TABLE IF NOT EXISTS chat.chat_messages_by_id (
    id         uuid PRIMARY KEY,
    room_id    uuid,
    created_at bigint
);
