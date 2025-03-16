-- +goose Up
CREATE TABLE IF NOT EXISTS PasswordResetRequest (
    id VARCHAR(255) NOT NULL,
    adminId VARCHAR(255) NOT NULL UNIQUE,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(adminId) REFERENCES Admin(id),
    INDEX(id, adminId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS PasswordResetRequest;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
