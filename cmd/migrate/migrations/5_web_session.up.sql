CREATE TABLE web_sessions (
    id SERIAL PRIMARY KEY,
    session_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    user_agent TEXT,
    ip_address VARCHAR(45),
    issued_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_web_sessions_account_id ON web_sessions(account_id);
CREATE INDEX idx_web_sessions_session_id ON web_sessions(session_id);
CREATE INDEX idx_web_sessions_expires_at ON web_sessions(expires_at);

CREATE TRIGGER update_web_sessions_updated_at 
    BEFORE UPDATE ON web_sessions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
