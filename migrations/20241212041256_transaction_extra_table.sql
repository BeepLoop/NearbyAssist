-- +goose Up
CREATE TABLE IF NOT EXISTS TransactionExtra (
    transactionId VARCHAR(255) NOT NULL,
    extraId VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY(transactionId) REFERENCES Transaction(id) ON DELETE CASCADE,
    FOREIGN KEY(extraId) REFERENCES Extra(id) ON DELETE CASCADE,
    INDEX(transactionId, extraId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS TransactionExtra;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
