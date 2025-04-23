-- +goose Up
CREATE TABLE IF NOT EXISTS ServiceAddress (
    serviceId VARCHAR(255) NOT NULL,
    addressId VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(serviceId) REFERENCES Service(id),
    FOREIGN KEY(addressId) REFERENCES Address(id),
    INDEX(serviceId, addressId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS ServiceAddress;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
