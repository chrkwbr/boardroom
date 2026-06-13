-- 000001_initial.cql
CREATE TABLE IF NOT EXISTS chat.chat_messages (
    room_id    uuid,
    id         uuid,
    sender_id  uuid,
    message    text,
    version    bigint,
    created_at bigint,
    updated_at bigint,
    PRIMARY KEY ((room_id), id)
) WITH CLUSTERING ORDER BY (id ASC);

CREATE TABLE IF NOT EXISTS chat.chat_message_histories (
    id         uuid,
    sender_id  uuid,
    message    text,
    version    bigint,
    room_id    uuid,
    status     text,
    created_at bigint,
    updated_at bigint,
    PRIMARY KEY ((room_id, id), version)
) WITH CLUSTERING ORDER BY (version DESC);
