-- should be executed by boardroom
create schema chat;

set search_path to chat;

CREATE TABLE IF NOT EXISTS chat_events (
    id SERIAL PRIMARY KEY,
    chat_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    version INTEGER NOT NULL,
    payload JSONB NOT NULL,
    created_at BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_chat_events_chat_id ON chat_events (chat_id);
ALTER TABLE chat_events ADD CONSTRAINT unique_chat_id_version UNIQUE (chat_id, version);

CREATE TABLE IF NOT EXISTS chat_outbox (
    id SERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL UNIQUE,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    processed BOOLEAN DEFAULT false,
    created_at BIGINT NOT NULL
);