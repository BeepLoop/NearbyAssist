-- +goose Up
CREATE TABLE IF NOT EXISTS Transaction (
    id VARCHAR(255) NOT NULL,
    vendorId VARCHAR(255) NOT NULL,
    clientId VARCHAR(255) NOT NULL,
    serviceId VARCHAR(255) NOT NULL,
    status Enum('pending', 'confirmed', 'rejected', 'done', 'cancelled') NOT NULL DEFAULT 'pending',
    cost DOUBLE NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES User(id) ON DELETE CASCADE,
    FOREIGN KEY(serviceId) REFERENCES Service(id) ON DELETE CASCADE,
    FOREIGN KEY(clientId) REFERENCES User(id),
    INDEX(id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Transaction;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
