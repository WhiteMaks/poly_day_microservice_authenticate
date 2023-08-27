CREATE TABLE root.users(
    id              BIGSERIAL       NOT NULL PRIMARY KEY,
    username        VARCHAR(25)     NOT NULL,
    password        VARCHAR(256)    NOT NULL,
    role_id         BIGSERIAL       NOT NULL,
    created_at      timestamptz     NOT NULL DEFAULT now()
);

ALTER TABLE root.users ADD CONSTRAINT FK_user_role_id FOREIGN KEY (role_id) REFERENCES root.users (id);