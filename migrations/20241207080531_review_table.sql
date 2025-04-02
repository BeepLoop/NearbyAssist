-- +goose Up
CREATE TABLE IF NOT EXISTS Review (
    id VARCHAR(255) NOT NULL,
    transactionId VARCHAR(255) NOT NULL,
    rating INT NOT NULL,
    text TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(transactionId) REFERENCES Transaction(id) ON DELETE CASCADE,
    INDEX(id, transactionId)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Review;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
