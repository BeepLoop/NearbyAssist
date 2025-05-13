-- +goose Up
CREATE TABLE IF NOT EXISTS Booking (
    id VARCHAR(255) NOT NULL,
    vendorId VARCHAR(255) NOT NULL,
    clientId VARCHAR(255) NOT NULL,
    serviceId VARCHAR(255) NOT NULL,
    serviceTitle VARCHAR(255) NOT NULL,
    serviceDescription TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    pricingType ENUM('fixed', 'per_hour', 'per_day') NOT NULL DEFAULT 'fixed',
    status Enum('pending', 'confirmed', 'rejected', 'done', 'cancelled') NOT NULL DEFAULT 'pending',
    quantity INT NOT NULL DEFAULT 1, -- for fixed price service
    cost DECIMAL(10, 2) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    scheduleStart TIMESTAMP,
    scheduleEnd TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    cancelledBy VARCHAR(255),
    cancelReason TEXT,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES User(id) ON DELETE CASCADE,
    FOREIGN KEY(clientId) REFERENCES User(id),
    FOREIGN KEY(serviceId) REFERENCES Service(id),
    FOREIGN KEY(cancelledBy) REFERENCES User(id),
    INDEX(id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Booking;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
