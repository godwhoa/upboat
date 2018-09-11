-- Just a holding ground for schema ideas
-- Real stuff in pkg/postgres/migrations
set timezone TO 'UTC';

DROP TABLE IF EXISTS users, posts, comments;

CREATE TABLE users(
    id serial PRIMARY KEY,
    uid TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hash TEXT NOT NULL,
    created TIMESTAMP DEFAULT now(),
    deleted TIMESTAMP NULL
);

CREATE TABLE posts(
    id serial PRIMARY KEY,
    author_id INTEGER REFERENCES users(id),
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    created TIMESTAMP DEFAULT now(),
    deleted TIMESTAMP NULL
);

CREATE TABLE post_votes(
    id serial PRIMARY KEY,
    voter_id INTEGER REFERENCES users(id),
    post_id INTEGER REFERENCES posts(id),
    UNIQUE(voter_id, post_id),
    delta INTEGER NOT NULL CHECK(delta IN (-1, +1))
);

CREATE TABLE comments(
    id serial PRIMARY KEY,
    post_id INTEGER REFERENCES posts(id),
    parent_id INTEGER REFERENCES comments(id),
    commenter_id INTEGER REFERENCES users(id),
    body TEXT NOT NULL,
    created TIMESTAMP DEFAULT now(),
    deleted TIMESTAMP DEFAULT NULL
);

CREATE TABLE comment_votes(
    id serial PRIMARY KEY,
    voter_id INTEGER REFERENCES users(id),
    comment_id INTEGER REFERENCES comments(id),
    delta INTEGER NOT NULL CHECK(delta IN (-1, +1))
);

-- Init users
INSERT INTO users(username) VALUES
("pacninja"),
("pakku");

INSERT INTO users(uid, username, email, hash) VALUES 
('fa5521311ae25d54a087841a8cc8b8d8dba577e1d726a77af404d4c8c4c52191', 'godwhoa', 'whoa@a.com', '$2y$12$4yxt4RQ.l.4zCOAyv6mkzuukmiprfb6iCHSeQMOQaUlOvS4F8ekLS'),
('ae47f7d43d6200b756cfa2f9f4ae95923406f23e9c004f3d26a438c7536316c9', 'pacninja', 'ninja@a.com', '$2y$12$SYWiI28ymT45yNIR2ljk/eA79ZYwloD/d61GlyQbwnx2RjsZwAX/u'),
;

-- Make post
INSERT INTO posts(author_id, body) VALUES (1, 'I am godwhoa');
INSERT INTO posts(author_id, title, body) VALUES
(1, 'Hello', 'Im pac'),
(2, 'Hello as well', 'Im dan');
-- Make comments
INSERT INTO comments(post_id, parent_id, commenter_id, body) VALUES
(1, NULL, 2, 'I am pacninja'),
(1, 1, 1, 'Nice to meet you pac.');

-- Update if it already
INSERT INTO post_votes(voter_id, post_id, delta) VALUES
(1, 1, +1),
(1, 1, -1),
(2, 1, -1),
(3, 1, +1),
(3, 1, -1);

-- Sum of last vote by each user on a post
SELECT COALESCE(SUM(delta), 0) FROM (
  SELECT DISTINCT ON(voter_id) voter_id, delta
  FROM post_votes
  WHERE post_id = 1
  GROUP BY voter_id, created, delta
  ORDER BY voter_id, created DESC
) as last_votes;