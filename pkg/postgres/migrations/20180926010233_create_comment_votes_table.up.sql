CREATE TABLE comment_votes(
    id serial PRIMARY KEY,
    voter_id INTEGER REFERENCES users(id),
    comment_id INTEGER REFERENCES comments(id),
    delta INTEGER NOT NULL CHECK(delta IN (-1, +1))
);