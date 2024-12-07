-- +goose Up
CREATE TABLE IF NOT EXISTS ServicePhoto (
    id VARCHAR(255) NOT NULL,
    serviceId VARCHAR(255) NOT NULL,
    vendorId VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(serviceId) REFERENCES Service(id) ON DELETE CASCADE,
    FOREIGN KEY(vendorId) REFERENCES User(id) ON DELETE CASCADE,
    INDEX(id, serviceId, vendorId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS ServicePhoto;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
