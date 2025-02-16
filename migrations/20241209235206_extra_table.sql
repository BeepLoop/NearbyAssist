-- +goose Up
CREATE TABLE IF NOT EXISTS Extra (
    id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    serviceId VARCHAR(255) NOT NULL,
    deleted TINYINT(1) NOT NULL DEFAULT 0 COMMENT '0: not deleted, 1: deleted',
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (serviceId) REFERENCES Service(id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Extra;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
