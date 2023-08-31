CREATE TABLE IF NOT EXISTS sessions (
    user_id  UUID PRIMARY KEY,
    session_token VARCHAR,
    expired_at TIMESTAMP,
)