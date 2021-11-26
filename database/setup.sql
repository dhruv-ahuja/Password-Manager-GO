DROP TABLE IF EXISTS info;

CREATE TABLE info (
    id SERIAL PRIMARY KEY, 
    key VARCHAR(128) NOT NULL, 
    encrypted_pw VARCHAR(128) NOT NULL);
    -- CONSTRAINT either_field
    -- CHECK (email is not null or username is not null)
