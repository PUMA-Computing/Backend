CREATE TABLE IF NOT EXISTS puma.news (
                           id BIGINT PRIMARY KEY,
                           author_id UUID,
                           title TEXT,
                           content TEXT,
                           visible BOOLEAN,
                           thumbnail TEXT,
                           categories SET<TEXT>,
                           published_date TIMESTAMP,
);