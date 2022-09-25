CREATE TABLE snippet (
     id SERIAL NOT NULL PRIMARY KEY,
     title VARCHAR(50) NOT NULL,
     content TEXT NOT NULL,
     created TIMESTAMP NOT NULL DEFAULT now(),
     expires TIMESTAMP NOT NULL
);