-- +goose Up
CREATE TABLE IF NOT EXISTS ActivityLog (
    id VARCHAR(255) NOT NULL,
    adminId VARCHAR(255) NOT NULL,
    action VARCHAR(255) NOT NULL,
    targetType VARCHAR(255) NOT NULL,
    targetId VARCHAR(255),
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(adminId) REFERENCES Admin(id),
    INDEX(id, adminId, targetId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS ActivityLog;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
