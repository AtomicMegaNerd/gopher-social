-- NOTE: Migrations will probably fail when we are trying to add foreign keys
-- to existing tables. For this reason we will set the values to a default int
-- afterwards we can query the roles table to actually set the proper default
-- that we want and then get rid of the default of 1 so it will not apply to
-- new users. None of this is needed if the database is not seeded yet.
ALTER TABLE
	IF EXISTS users
ADD
	COLUMN role_id INT REFERENCES roles(id) DEFAULT 1;

UPDATE users
SET role_id = (
	SELECT id FROM roles WHERE name = 'user'
);

ALTER TABLE
	IF EXISTS users
ALTER COLUMN role_id
SET NOT NULL;

ALTER TABLE
	IF EXISTS users
ALTER COLUMN role_id
DROP DEFAULT;
