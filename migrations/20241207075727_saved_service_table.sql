-- +goose Up
CREATE TABLE IF NOT EXISTS ServiceTag (
    id VARCHAR(255) NOT NULL,
    serviceId VARCHAR(255) NOT NULL,
    tagId VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(serviceId) REFERENCES Service(id) ON DELETE CASCADE,
    FOREIGN KEY(tagId) REFERENCES Tag(id),
    INDEX(id, serviceId, tagId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS ServiceTag;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
