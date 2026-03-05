-- Table for tracking message status: sent, delivered, read (like WhatsApp ✓✓)
-- Each message has 1 receipt per recipient (sender doesn't need a receipt for themselves)
CREATE TABLE message_receipts (
    -- Unique ID for each receipt (rcpt_xxx)
    id VARCHAR(32) PRIMARY KEY,

    -- Which message is being tracked (FK to messages table)
    message_id VARCHAR(32) NOT NULL,

    -- Which recipient (FK to users table) — not the sender, but the receiver
    user_id VARCHAR(32) NOT NULL,

    -- Receipt status: 'sent', 'delivered', or 'read'
    status VARCHAR(20) NOT NULL DEFAULT 'sent',

    -- When the message arrived at recipient's device (null = not yet delivered)
    delivered_at TIMESTAMP,

    -- When the recipient read the message (null = not yet read)
    read_at TIMESTAMP,

    -- Foreign keys
    CONSTRAINT fk_receipt_message FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    CONSTRAINT fk_receipt_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    -- 1 message can only have 1 receipt per user (no duplicates)
    CONSTRAINT uq_message_user UNIQUE (message_id, user_id)
);

-- Index for fast query: "all receipts for message X"
CREATE INDEX idx_receipts_message_id ON message_receipts(message_id);

-- Index for fast query: "all receipts for user Y"
CREATE INDEX idx_receipts_user_id ON message_receipts(user_id);