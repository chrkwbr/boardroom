-- 000002_add_status_to_chat_message_histories.cql
ALTER TABLE chat.chat_message_histories
ADD status text;

