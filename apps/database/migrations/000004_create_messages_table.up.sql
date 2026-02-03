CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(32) PRIMARY KEY,
    conversation_id VARCHAR(32) NOT NULL,
    sender_id VARCHAR(32) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(20) DEFAULT 'text',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_messages_conversation 
        FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT fk_messages_sender 
        FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index untuk mengambil messages by conversation (most common query)
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);

-- Index untuk soft delete
CREATE INDEX idx_messages_deleted_at ON messages(deleted_at);

-- Index untuk sorting by time (descending untuk chat terbaru)
CREATE INDEX idx_messages_created_at ON messages(created_at DESC);