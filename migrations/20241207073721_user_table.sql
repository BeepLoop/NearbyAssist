-- +goose Up
CREATE TABLE IF NOT EXISTS User (
    id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    emailHash VARCHAR(64) NOT NULL,
    imageUrl VARCHAR(255),
    phone VARCHAR(255),
    verified BOOLEAN DEFAULT FALSE,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verifiedAt TIMESTAMP,
    PRIMARY KEY (id),
    INDEX(id, name, emailHash)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS User;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
