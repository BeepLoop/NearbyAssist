-- +goose Up
CREATE TABLE IF NOT EXISTS Expertise (
    id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Expertise;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
