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
