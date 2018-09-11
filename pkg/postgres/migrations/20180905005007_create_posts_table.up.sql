CREATE TABLE posts(
    id serial PRIMARY KEY,
    author_id INTEGER REFERENCES users(id),
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    created TIMESTAMP DEFAULT now(),
    deleted TIMESTAMP NULL
);