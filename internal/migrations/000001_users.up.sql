CREATE TABLE IF NOT EXISTS users (
                       id UUID PRIMARY KEY,
                       first_name VARCHAR(255),
                       last_name VARCHAR(255),
                       email VARCHAR(255) NOT NULL UNIQUE ,
                       password VARCHAR(255) NOT NULL,
                       role_id INTEGER REFERENCES roles(id),
                       nim VARCHAR(255) UNIQUE,
                       year VARCHAR(255),
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);