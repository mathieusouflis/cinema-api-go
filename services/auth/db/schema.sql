CREATE TABLE users (
    id        UUID        PRIMARY KEY,
    username  VARCHAR(15) UNIQUE NOT NULL,
    email     TEXT        UNIQUE NOT NULL,
    password  TEXT        NOT NULL,
    google_id TEXT        UNIQUE,
    github_id TEXT        UNIQUE,
    role      TEXT        NOT NULL DEFAULT 'user'
);
