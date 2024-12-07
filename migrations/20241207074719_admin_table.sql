-- +goose Up
CREATE TABLE IF NOT EXISTS Admin (
    id VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role Enum('admin', 'staff') NOT NULL DEFAULT 'staff',
    usernameHash VARCHAR(64) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(id, usernameHash, role)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Admin;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
