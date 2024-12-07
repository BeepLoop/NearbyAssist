-- +goose Up
CREATE TABLE IF NOT EXISTS Application (
    id VARCHAR(255) NOT NULL,
    applicantId VARCHAR(255) NOT NULL UNIQUE,
    job VARCHAR(255) NOT NULL,
    status Enum('pending', 'rejected', 'approved') NOT NULL DEFAULT 'pending',
    supportingDocumentUrl VARCHAR(255) NOT NULL,
    policeClearanceUrl VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(applicantId) REFERENCES User(id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Application;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
