-- +goose Up
CREATE TABLE IF NOT EXISTS Invitation (
    id CHAR(36) DEFAULT (UUID()),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    usernameHash VARCHAR(255) NOT NULL,
    emailHash VARCHAR(255) NOT NULL UNIQUE,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expiredAt TIMESTAMP NOT NULL,
    PRIMARY KEY(id),
    INDEX (id, email, code, expiredAt)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Invitation;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
