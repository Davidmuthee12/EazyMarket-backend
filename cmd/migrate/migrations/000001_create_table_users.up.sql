CREATE EXTENSION IF NOT EXISTS citext; 


CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    username varchar(255) UNIQUE NOT NULL,
    phone varchar(20),
    avatar_url text,
    role varchar(20) NOT NULL DEFAULT 'user', 
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    password bytea NOT NULL, 
    created_at timestamp(0) with time zone NOT NULL DEFAULt now()
)