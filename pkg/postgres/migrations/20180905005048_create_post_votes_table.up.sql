CREATE TABLE post_votes(
    id serial PRIMARY KEY,
    voter_id INTEGER REFERENCES users(id),
    post_id INTEGER REFERENCES posts(id),
    UNIQUE(voter_id, post_id),
    delta INTEGER NOT NULL CHECK(delta IN (-1, +1))
);