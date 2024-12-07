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
        WHERE serviceId = NEW.serviceId
    )
    WHERE vendorId = (
        SELECT vendorId 
        FROM Service 
        WHERE id = NEW.serviceId
    );
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
