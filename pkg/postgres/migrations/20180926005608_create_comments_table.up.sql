CREATE TABLE comments(
    id serial PRIMARY KEY,
    post_id INTEGER REFERENCES posts(id),
    parent_id INTEGER REFERENCES comments(id),
    commenter_id INTEGER REFERENCES users(id),
    depth INTEGER NOT NULL,
    body TEXT NOT NULL,
    created TIMESTAMP DEFAULT now(),
    deleted TIMESTAMP DEFAULT NULL
);