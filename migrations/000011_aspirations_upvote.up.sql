SET timezone = 'Asia/Jakarta';
CREATE TABLE IF NOT EXISTS aspirations_upvote (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    aspiration_id int NOT NULL REFERENCES aspirations(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, aspiration_id)
);