-- +goose Up
CREATE TABLE IF NOT EXISTS Review (
    id CHAR(36) DEFAULT (UUID()),
    bookingId VARCHAR(255) NOT NULL,
    rating INT NOT NULL,
    text TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(bookingId) REFERENCES Booking(id) ON DELETE CASCADE,
    INDEX(id, bookingId)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS Review;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
