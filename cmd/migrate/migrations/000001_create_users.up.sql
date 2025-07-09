CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL, -- using citext for case-insensitive email
    username varchar(255) UNIQUE NOT NULL,
    password bytea NOT NULL, -- storing hashed password
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);
