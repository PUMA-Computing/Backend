SET timezone = 'Asia/Jakarta';
CREATE TABLE IF NOT EXISTS event_registrations (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL,
    user_id uuid NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
)