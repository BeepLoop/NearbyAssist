-- +goose Up
CREATE TABLE IF NOT EXISTS User (
    id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    emailHash VARCHAR(64) NOT NULL,
    imageUrl VARCHAR(255),
    address VARCHAR(255),
    phone VARCHAR(255),
    latitude Decimal(12, 10),
    longitude Decimal(13, 10),
    verified TINYINT(1) NOT NULL DEFAULT 0 COMMENT '0: not verified, 1: verified',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
