-- CREATE TABLE users
-- (
--     id serial not null unique,
--     name varchar(255) not null,
--     balance NUMERIC(20, 2) DEFAULT 0 NOT NULL,
--     username varchar(255) not null unique,
--     password varchar(255) not null
-- );
CREATE TABLE payment_sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);