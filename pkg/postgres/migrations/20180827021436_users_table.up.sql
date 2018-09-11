CREATE TABLE users(
    id serial PRIMARY KEY,
    uid TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hash TEXT NOT NULL,
    created TIMESTAMP DEFAULT now(),
    deleted TIMESTAMP NULL
);