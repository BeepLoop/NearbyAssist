-- +goose Up
CREATE TABLE IF NOT EXISTS Restricted (
    userId VARCHAR(255) NOT NULL UNIQUE,
    reason TEXT,
    startTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    endTime TIMESTAMP NOT NULL,
    FOREIGN KEY(userId) REFERENCES User(id),
    INDEX(userId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Restricted;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
