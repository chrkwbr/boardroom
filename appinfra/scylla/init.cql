CREATE KEYSPACE IF NOT EXISTS chat
    WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE IF NOT EXISTS chat.chat_messages (
    room_id     uuid,
    id          uuid,
    message     text,
    sender_id   uuid,
    version     bigint,
    created_at  bigint,
    updated_at  bigint,
    PRIMARY KEY ((room_id), id)
) WITH CLUSTERING ORDER BY (id DESC);

CREATE TABLE IF NOT EXISTS chat.chat_message_histories (
    id          uuid,
    message     text,
    version     bigint,
    room_id     uuid,
    sender_id   uuid,
    created_at  bigint,
    updated_at  bigint,
    PRIMARY KEY ((id), version)
) WITH CLUSTERING ORDER BY (version DESC);
