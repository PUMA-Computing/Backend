CREATE TABLE IF NOT EXISTS events (
                                      id SERIAL PRIMARY KEY,
                                      title TEXT NOT NULL,
                                      description TEXT NOT NULL,
                                      date DATE,
                                      user_id INTEGER NOT NULL,
                                      created_at TIMESTAMP DEFAULT NOW(),
                                      updated_at TIMESTAMP DEFAULT NOW(),
                                      FOREIGN KEY (user_id) REFERENCES users (id)
);