-- +goose Up
CREATE TABLE IF NOT EXISTS BookingExtra (
    bookingId VARCHAR(255) NOT NULL,
    extraTitle VARCHAR(255) NOT NULL,
    extraDescription TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY(bookingId) REFERENCES Booking(id) ON DELETE CASCADE,
    INDEX(bookingId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS BookingExtra;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
