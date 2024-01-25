CREATE TABLE IF NOT EXISTS event_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
    );

INSERT INTO event_categories (name)
VALUES
    ('PUFA Computing'),
    ('PUMA IT'),
    ('PUMA IS'),
    ('PUMA ID'),
    ('PUMA VCD');

ALTER TABLE events
    ADD COLUMN category_id INT NOT NULL,
    ADD FOREIGN KEY (category_id) REFERENCES event_categories(id);
