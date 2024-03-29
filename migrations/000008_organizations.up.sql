SET timezone = 'Asia/Jakarta';
CREATE TABLE IF NOT EXISTS organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
    );

INSERT INTO organizations (name)
VALUES
    ('PUFA Computing'),
    ('PUMA IT'),
    ('PUMA IS'),
    ('PUMA ID'),
    ('PUMA VCD');

ALTER TABLE events
    ADD COLUMN organization_id INT NOT NULL,
    ADD FOREIGN KEY (organization_id) REFERENCES organizations(id);
