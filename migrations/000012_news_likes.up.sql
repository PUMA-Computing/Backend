SET timezone = 'Asia/Jakarta';
CREATE TABLE IF NOT EXISTS news_likes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    news_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (news_id) REFERENCES news(id),
    UNIQUE (user_id, news_id)
)