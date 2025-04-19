-- +goose Up
CREATE TABLE IF NOT EXISTS Address (
    id VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    latitude Decimal(12, 10) NOT NULL,
    longitude Decimal(13, 10) NOT NULL,
    PRIMARY KEY(id),
    INDEX(id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Address;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
