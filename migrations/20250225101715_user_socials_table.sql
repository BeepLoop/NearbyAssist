-- +goose Up
CREATE TABLE IF NOT EXISTS Social(
    id VARCHAR(255) NOT NULL,
    userId VARCHAR(255) NOT NULL,
    site VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(userId) REFERENCES User(id),
    INDEX(id, userId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Social;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
