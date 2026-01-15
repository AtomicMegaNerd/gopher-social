CREATE TABLE IF NOT EXISTS roles (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE,
	level int NOT NULL DEFAULT 0,
	description TEXT
);

-- Add the core roles with the migration
INSERT INTO roles (name, level, description) VALUES
('user', 1, 'A user can create posts and comments'),
('moderator', 2, 'A moderator can manage user content'),
('admin', 3, 'An admin has full access to the system');
