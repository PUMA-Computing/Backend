CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO permissions (name, description)
VALUES
    ('get:users', 'Get users by id'),
    ('create:users', 'Create users'),
    ('edit:users', 'Edit users'),
    ('delete:users', 'Delete users'),
    ('list:users', 'List users'),
    ('get:events', 'Get events by id'),
    ('create:events', 'Create events'),
    ('edit:events', 'Edit events'),
    ('delete:events', 'Delete events'),
    ('list:events', 'List events'),
    ('register:events', 'Register for events'),
    ('listRegisterEvent:events', 'List user registered for event'),
    ('get:news', 'Get news'),
    ('create:news', 'Create news'),
    ('edit:news', 'Edit news'),
    ('delete:news', 'Delete news'),
    ('list:news', 'List news'),
    ('like:news', 'Like news'),
    ('get:roles', 'Get roles'),
    ('create:roles', 'Create roles'),
    ('edit:roles', 'Edit roles'),
    ('delete:roles', 'Delete roles'),
    ('list:roles', 'List roles'),
    ('assign:roles', 'Assign Role to User'),
    ('list:permissions', 'List permissions'),
    ('assign:permissions', 'Assign permissions to role');