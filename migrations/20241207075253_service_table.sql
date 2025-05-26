-- +goose Up
CREATE TABLE IF NOT EXISTS Service (
    id VARCHAR(255) NOT NULL,
    vendorId VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    pricingType ENUM('fixed', 'per_hour', 'per_day') NOT NULL DEFAULT 'fixed',
    signature VARCHAR(64) NOT NULL,
    disabled BOOLEAN DEFAULT FALSE,
    status ENUM('under_review', 'accepted', 'rejected') DEFAULT 'under_review',
    rejectReason TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    acceptedAt TIMESTAMP,
    rejectedAt TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES Vendor(vendorId),
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
