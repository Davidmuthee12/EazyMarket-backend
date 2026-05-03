CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO 
    roles (name, description, level)
VALUES
    (
        'user',
        'A user can create posts and comment',
        1
    );

INSERT INTO 
    roles (name, description, level)
VALUES
    (
        'vendor',
        'A vendor can sell products in the platform',
        2
    );

INSERT INTO
    roles (name, description, level)
VALUES
    (
        'admin',
        'An Admin is the general maintainer and owner of the platform',
        3
    );