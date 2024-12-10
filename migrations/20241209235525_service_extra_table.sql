-- +goose Up
CREATE TABLE IF NOT EXISTS ServiceExtra (
    serviceId VARCHAR(255) NOT NULL,
    extraId VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY(serviceId) REFERENCES Service(id),
    FOREIGN KEY(extraId) REFERENCES Extra(id),
    INDEX(serviceId, extraId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS ServiceExtra;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
