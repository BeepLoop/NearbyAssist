-- +goose Up
CREATE TABLE IF NOT EXISTS IdentityVerification (
    id VARCHAR(255) NOT NULL,
    userId VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    latitude Decimal(12, 10) NOT NULL,
    longitude Decimal(13, 10) NOT NULL,
    idType VARCHAR(255) NOT NULL,
    idNumber VARCHAR(255) NOT NULL,
    frontIdImageUrl VARCHAR(255) NOT NULL,
    backIdImageUrl VARCHAR(255) NOT NULL,
    faceImageUrl VARCHAR(255) NOT NULL,
    status Enum('pending', 'rejected', 'approved') NOT NULL DEFAULT 'pending',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id, userId, idType),
    FOREIGN KEY(userId) REFERENCES User(id),
    INDEX(id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS IdentityVerification;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
