-- +goose Up
CREATE TABLE IF NOT EXISTS Identification (
    id VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    referenceNumber VARCHAR(255) NOT NULL,
    frontImageUrl VARCHAR(255) NOT NULL,
    backImageUrl VARCHAR(255) NOT NULL,
    selfieImageUrl VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Identification;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
