-- +goose Up
CREATE TABLE IF NOT EXISTS IdentityVerification (
    id VARCHAR(255) NOT NULL,
    userId VARCHAR(255) NOT NULL,
    status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
    rejectionNote TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(userId) REFERENCES User(id),
    INDEX(id, userId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS IdentityVerification;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
