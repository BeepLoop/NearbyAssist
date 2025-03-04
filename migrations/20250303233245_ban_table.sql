-- +goose Up
CREATE TABLE IF NOT EXISTS Ban (
    userId VARCHAR(255) NOT NULL UNIQUE,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(userId) REFERENCES User(id),
    INDEX (userId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Ban;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
