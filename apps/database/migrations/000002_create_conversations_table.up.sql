CREATE TABLE IF NOT EXISTS conversations (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(100),
    type VARCHAR(20) NOT NULL DEFAULT 'direct',
    created_by VARCHAR(32) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Index untuk soft delete
CREATE INDEX idx_conversations_deleted_at ON conversations(deleted_at);

-- Index untuk creator
CREATE INDEX idx_conversations_created_by ON conversations(created_by);