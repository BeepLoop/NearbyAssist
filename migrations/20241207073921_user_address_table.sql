-- +goose Up
CREATE TABLE IF NOT EXISTS UserAddress (
    userId VARCHAR(255) NOT NULL,
    addressId VARCHAR(36) NOT NULL,
    FOREIGN KEY(userId) REFERENCES User(id),
    FOREIGN KEY(addressId) REFERENCES Address(id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS UserAddress;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
