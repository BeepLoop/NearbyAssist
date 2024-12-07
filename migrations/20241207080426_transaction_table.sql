-- +goose Up
CREATE TABLE IF NOT EXISTS Transaction (
    id VARCHAR(255) NOT NULL,
    vendorId VARCHAR(255) NOT NULL,
    clientId VARCHAR(255) NOT NULL,
    serviceId VARCHAR(255) NOT NULL,
    status Enum('pending', 'ongoing', 'done', 'cancelled') NOT NULL DEFAULT 'pending',
    employmentType Enum('arawan', 'pakyaw') NOT NULL DEFAULT 'pakyaw',
    cost DOUBLE NOT NULL,
    startDate TIMESTAMP NOT NULL,
    endDate TIMESTAMP NOT NULL,
    isReviewed TINYINT(1) NOT NULL DEFAULT 0 COMMENT '0: not reviewed, 1: reviewed',
    isReported TINYINT(1) NOT NULL DEFAULT 0 COMMENT '0: not reported, 1: reported',
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
