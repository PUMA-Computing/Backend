ALTER TABLE IF EXISTS users
ADD COLUMN IF NOT EXISTS "twofa_enabled" BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS "twofa_image" TEXT,
ADD COLUMN IF NOT EXISTS "twofa_secret" TEXT;

INSERT INTO permissions (name, description)
VALUES
    ('users:2fa', '2FA Feature');

INSERT INTO role_permissions (role_id, permission_id)
VALUES
    (1, 34),
    (2, 34),
    (3, 34),
    (4, 34),
    (5, 34),
    (6, 34),
    (7, 34),
    (8, 34);