CREATE TABLE snippet (
     id SERIAL NOT NULL PRIMARY KEY,
     title VARCHAR(50) NOT NULL,
     content TEXT NOT NULL,
     created TIMESTAMP NOT NULL DEFAULT now(),
     expires TIMESTAMP NOT NULL
);

CREATE TABLE sessions (
    token VARCHAR(43) PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMP  NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);