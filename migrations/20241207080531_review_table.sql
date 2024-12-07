-- +goose Up
CREATE TABLE IF NOT EXISTS Review (
    id VARCHAR(255) NOT NULL,
    serviceId VARCHAR(255) NOT NULL,
    rating INT NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(serviceId) REFERENCES Service(id) ON DELETE CASCADE,
    INDEX(id, serviceId)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Review;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
