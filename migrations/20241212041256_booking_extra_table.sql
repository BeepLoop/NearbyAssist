-- +goose Up
CREATE TABLE IF NOT EXISTS BookingExtra (
    bookingId VARCHAR(255) NOT NULL,
    extraId VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY(bookingId) REFERENCES Booking(id) ON DELETE CASCADE,
    FOREIGN KEY(extraId) REFERENCES Extra(id) ON DELETE CASCADE,
    INDEX(bookingId, extraId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS BookingExtra;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
