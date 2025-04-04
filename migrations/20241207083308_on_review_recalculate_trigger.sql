-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER update_vendor_rating
AFTER INSERT ON Review
FOR EACH ROW
BEGIN
    -- Update the average rating for the vendor based on the new review
    UPDATE Vendor
    SET rating = (
        SELECT ROUND(AVG(rating), 1)
        FROM Review
        WHERE bookingId = NEW.bookingId
    )
    WHERE vendorId = (
        SELECT vendorId 
        FROM Booking
        WHERE id = NEW.bookingId
    );
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS update_vendor_rating;
-- +goose StatementBegin
-- +goose StatementEnd
