-- +goose Up
CREATE TABLE IF NOT EXISTS Vendor (
    vendorId VARCHAR(255) NOT NULL,
    dbl INT NOT NULL DEFAULT 5,
    rating Decimal(5,1) NOT NULL DEFAULT 0.0,
    joinedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(vendorId),
    FOREIGN KEY(vendorId) REFERENCES User(id),
    INDEX(vendorId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Vendor;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
