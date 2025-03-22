-- +goose Up
CREATE TABLE IF NOT EXISTS SupportingImage (
    id VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX (id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS SupportingImage;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
