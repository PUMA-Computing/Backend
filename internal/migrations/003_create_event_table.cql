CREATE TABLE IF NOT EXISTS events (
                                     id BIGINT PRIMARY KEY,
                                     image TEXT,
                                     name TEXT,
                                     description TEXT,
                                     date TIMESTAMP,
                                     registered_users SET<TEXT>
);