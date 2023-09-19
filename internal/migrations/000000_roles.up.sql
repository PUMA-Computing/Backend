CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    prefix VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id) REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS role_permissions (
    id SERIAL PRIMARY KEY,
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

INSERT INTO roles (name)
VALUES
    ('PUMA'),
    ('Computizen');

INSERT INTO permissions (name, prefix)
VALUES
    ('Manage All', 'api_admin_*'),
    ('Manage Profile', 'api_user_*'),
    ('Edit Profile', 'api_user_editProfile'),
    ('Change Password', 'api_user_changePassword'),
    ('Delete Account', 'api_user_deleteAccount'),
    ('Register Event', 'api_user_registerEvent'),
    ('Manage Users', 'api_admin_manage_users_*'),
    ('Create User', 'api_admin_manage_users_createUser'),
    ('Edit User', 'api_admin_manage_users_editUser'),
    ('Delete User', 'api_admin_manage_users_deleteUser'),
    ('Manage Event', 'api_admin_manage_events_*'),
    ('Create Event', 'api_admin_manage_events_createEvent'),
    ('Edit Event', 'api_admin_manage_events_editEvent'),
    ('Delete Event', 'api_admin_manage_events_deleteEvent'),
    ('Get Users Registered Event', 'api_admin_manage_events_getUsersRegisteredEvent'),
    ('Manage Users Registered Event', 'api_admin_manage_events_manageUsersRegisteredEvent'),
    ('Delete Users Registered Event', 'api_admin_manage_events_deleteUsersRegisteredEvent'),
    ('Manage News', 'api_admin_manage_news_*'),
    ('Create News', 'api_admin_manage_news_createNews'),
    ('Edit News', 'api_admin_manage_news_editNews'),
    ('Delete News', 'api_admin_manage_news_deleteNews'),
    ('Manage Roles & Permissions', 'api_admin_manage_roles_*'),
    ('Create Role', 'api_admin_manage_roles_createRole'),
    ('Edit Role', 'api_admin_manage_roles_editRole'),
    ('Delete Role', 'api_admin_manage_roles_deleteRole'),
    ('Assign Role Permissions', 'api_admin_manage_roles_assignRolePermissions'),
    ('Assign User Role', 'api_admin_manage_roles_assignUserRole');

INSERT INTO role_permissions (role_id, permission_id)
VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (1, 4),
    (1, 5),
    (1, 6),
    (1, 7),
    (1, 8),
    (1, 9),
    (1, 10),
    (1, 11),
    (1, 12),
    (1, 13),
    (1, 14),
    (1, 15),
    (1, 16),
    (1, 17),
    (1, 18),
    (1, 19),
    (1, 20),
    (1, 21),
    (1, 22),
    (1, 23),
    (1, 24),
    (1, 25),
    (1, 26),
    (2, 2),
    (2, 3),
    (2, 4),
    (2, 5),
    (2, 6);
