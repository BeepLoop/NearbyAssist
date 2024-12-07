-- +goose Up
CREATE TABLE IF NOT EXISTS ExpertiseTag (
    expertiseId VARCHAR(255) NOT NULL,
    tagId VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY(expertiseId) REFERENCES Expertise(id),
    FOREIGN KEY(tagId) REFERENCES Tag(id),
    INDEX(expertiseId, tagId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS ExpertiseTag;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
