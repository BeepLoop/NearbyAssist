-- +goose Up
CREATE TABLE IF NOT EXISTS Notification (
    id VARCHAR(255) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    type Enum('success', 'fail', 'announcement', 'generic') NOT NULL DEFAULT 'generic',
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    readAt TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(recipient) REFERENCES User(id) ON DELETE CASCADE,
    INDEX(id, recipient)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Notification;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
