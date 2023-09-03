CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    author_id uuid REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content VARCHAR NOT NULL,
    category_id INTEGER REFERENCES news_categories(id),
    thumbnail VARCHAR(255),
    status VARCHAR(255) DEFAULT 'draft',
    published_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)