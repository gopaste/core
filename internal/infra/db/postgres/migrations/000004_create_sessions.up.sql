CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    refresh_token TEXT,
    user_agent VARCHAR(255),
    client_ip VARCHAR(15),
    is_blocked BOOLEAN,
    is_revoked BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP
);
