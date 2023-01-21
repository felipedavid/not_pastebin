CREATE TABLE snippets (
	id SERIAL NOT NULL,
	title VARCHAR(100) NOT NULL,
	content TEXT NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	expires TIMESTAMP NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

INSERT INTO snippets (title, content, expires) VALUES (
	'An old silent pond',
	'An old silent pond...\nA frog jumps into the pond,\nsplash!',
	NOW() + INTERVAL '1 day');

INSERT INTO snippets (title, content, expires) VALUES (
	'Over the wintry forest',
	'Over the wintry\nforest, winds howl in rage.',
	NOW() + INTERVAL '365 day');

CREATE TABLE sessions (
    token VARCHAR(43) PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMP NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

CREATE TABLE users (
    id SERIAL NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD CONSTRAINT users_email_uc UNIQUE (email);
