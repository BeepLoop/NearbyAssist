-- +goose Up
CREATE TABLE IF NOT EXISTS Session (
    id VARCHAR(255) NOT NULL,
    refreshToken VARCHAR(255) NOT NULL,
    status Enum('online', 'offline') NOT NULL DEFAULT 'online',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    CONSTRAINT unique_refreshToken_online UNIQUE (refreshToken, (CASE WHEN status = 'online' THEN 1 ELSE NULL END)),
    INDEX(refreshToken, status)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Session;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
