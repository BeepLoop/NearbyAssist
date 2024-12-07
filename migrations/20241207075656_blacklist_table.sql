-- +goose Up
CREATE TABLE IF NOT EXISTS Blacklist (
    id VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(token)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Blacklist;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
