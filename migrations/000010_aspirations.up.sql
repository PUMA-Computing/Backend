CREATE TABLE IF NOT EXISTS aspirations (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    subject varchar(50) NOT NULL,
    message TEXT NOT NULL,
    anonymous boolean NOT NULL DEFAULT FALSE,
    organization_id int NOT NULL REFERENCES organizations(id),
    closed boolean NOT NULL DEFAULT FALSE,
    admin_reply TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);