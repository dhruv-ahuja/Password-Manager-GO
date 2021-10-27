DROP TABLE IF EXISTS info;

CREATE TABLE info (
    id SERIAL PRIMARY KEY, 
    website VARCHAR(64) NOT NULL,
    email VARCHAR(128),
    username VARCHAR(64), 
    password_hash VARCHAR(64) NOT NULL,
    CONSTRAINT either_field
    CHECK (email is not null or username is not null)
);