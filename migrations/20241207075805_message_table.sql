-- +goose Up
CREATE TABLE IF NOT EXISTS Message (
    id VARCHAR(255) NOT NULL,
    sender VARCHAR(255) NOT NULL,
    receiver VARCHAR(255) NOT NULL,
    content Text NOT NULL,
    seen BOOLEAN DEFAULT FALSE,
    seenAt TIMESTAMP,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(sender) REFERENCES User(id),
    FOREIGN KEY(receiver) REFERENCES User(id),
    INDEX(id, sender, receiver)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Message;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
