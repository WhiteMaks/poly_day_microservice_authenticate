CREATE TABLE root.roles(
    id      BIGSERIAL       NOT NULL    PRIMARY KEY,
    name    VARCHAR(25)     NOT NULL    UNIQUE
);

INSERT INTO
    root.roles(name)
VALUES
    ('ROLE_DEVELOPER'),
    ('ROLE_GAMER');