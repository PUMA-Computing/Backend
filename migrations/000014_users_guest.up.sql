ALTER TABLE users
ADD COLUMN institution_name VARCHAR(255),

ALTER COLUMN student_id DROP NOT NULL,
ALTER COLUMN student_id SET DEFAULT NULL,
ALTER COLUMN major DROP NOT NULL,
ALTER COLUMN major SET DEFAULT NULL;

INSERT INTO roles (name) VALUES ('guest');
INSERT INTO role_permissions (role_id, permission_id)
VALUES
    (8, 1),
    (8, 6),
    (8, 11),
    (8, 13),
    (8, 18),
    (8, 2),
    (8, 3),
    (8, 27);

ALTER TABLE events
ADD COLUMN open_for_all BOOLEAN NOT NULL DEFAULT FALSE;