-- +goose Up
CREATE TABLE IF NOT EXISTS SavedService (
    id VARCHAR(255) NOT NULL,
    userId VARCHAR(255) NOT NULL,
    serviceId VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(userId) REFERENCES User(id),
    FOREIGN KEY(serviceId) REFERENCES Service(id),
    INDEX(id, userId, serviceId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS SavedService;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
