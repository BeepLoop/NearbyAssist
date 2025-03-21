-- +goose Up
CREATE TABLE IF NOT EXISTS Vendor (
    vendorId VARCHAR(255) NOT NULL UNIQUE,
    rating Decimal(5,1) NOT NULL DEFAULT 0.0,
    joinedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES User(id),
    INDEX(id, vendorId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Vendor;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
