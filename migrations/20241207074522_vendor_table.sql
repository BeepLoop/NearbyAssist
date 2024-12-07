-- +goose Up
CREATE TABLE IF NOT EXISTS Vendor (
    id VARCHAR(255) NOT NULL,
    vendorId VARCHAR(255) NOT NULL,
    rating Decimal(5,1) NOT NULL DEFAULT 0.0,
    restricted TINYINT(1) NOT NULL DEFAULT 0,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
