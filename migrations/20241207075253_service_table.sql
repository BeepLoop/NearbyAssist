-- +goose Up
CREATE TABLE IF NOT EXISTS Service (
    id VARCHAR(255) NOT NULL,
    vendorId VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    rate Double NOT NULL,
    latitude Decimal(12, 10) NOT NULL,
    longitude Decimal(13, 10) NOT NULL,
    signature VARCHAR(64) NOT NULL,
    disabled BOOLEAN DEFAULT FALSE,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES User(id),
    INDEX(id, vendorId, signature)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Service;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
