-- 000006_add_is_edited_to_messages.up.sql
ALTER TABLE messages ADD COLUMN is_edited BOOLEAN NOT NULL DEFAULT FALSE;