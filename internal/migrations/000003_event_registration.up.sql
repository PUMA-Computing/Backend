CREATE TABLE IF NOT EXISTS event_registration (
    id SERIAL PRIMARY KEY,
    event_id INTEGER REFERENCES events(id),
    user_id uuid REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)