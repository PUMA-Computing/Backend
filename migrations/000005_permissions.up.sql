CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO permissions (name, description)
VALUES
    ('users:get', 'Get users by id'),
    ('users:create', 'Create users'),
    ('users:edit', 'Edit users'),
    ('users:delete', 'Delete users'),
    ('users:list', 'List users'),
    ('events:get', 'Get events by id'),
    ('events:create', 'Create events'),
    ('events:edit', 'Edit events'),
    ('events:delete', 'Delete events'),
    ('events:list', 'List events'),
    ('events:register', 'Register for events'),
    ('events:listRegisteredUsers', 'List user registered for event'),
    ('news:get', 'Get news'),
    ('news:create', 'Create news'),
    ('news:edit', 'Edit news'),
    ('news:delete', 'Delete news'),
    ('news:list', 'List news'),
    ('news:like', 'Like news'),
    ('roles:get', 'Get roles'),
    ('roles:create', 'Create roles'),
    ('roles:edit', 'Edit roles'),
    ('roles:delete', 'Delete roles'),
    ('roles:list', 'List roles'),
    ('roles:assign', 'Assign Role to User'),
    ('permissions:list', 'List permissions'),
    ('permissions:assign', 'Assign permissions to role');