CREATE TABLE IF NOT EXISTS participants (
    id VARCHAR(32) PRIMARY KEY,
    user_id VARCHAR(32) NOT NULL,
    conversation_id VARCHAR(32) NOT NULL,
    role VARCHAR(20) DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_participants_user 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_participants_conversation 
        FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    
    -- Unique constraint: 1 user hanya bisa 1x join per conversation
    CONSTRAINT unique_user_conversation UNIQUE (user_id, conversation_id)
);

-- Index untuk mencari conversations by user
CREATE INDEX idx_participants_user_id ON participants(user_id);

-- Index untuk mencari users by conversation
CREATE INDEX idx_participants_conversation_id ON participants(conversation_id);