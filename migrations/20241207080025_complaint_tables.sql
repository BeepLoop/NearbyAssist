-- +goose Up
CREATE TABLE IF NOT EXISTS VendorComplaint (
    id VARCHAR(255) NOT NULL,
    vendorId VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content Text NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES User(id) ON DELETE CASCADE,
    INDEX(id, vendorId)
);

CREATE TABLE IF NOT EXISTS SystemComplaint (
    id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    detail Text NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(id, title)
);

CREATE TABLE IF NOT EXISTS SystemComplaintImage (
    id VARCHAR(255) NOT NULL,
    complaintId VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(complaintId) REFERENCES SystemComplaint(id),
    INDEX(id, complaintId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS SystemComplaintImage;
DROP TABLE IF EXISTS SystemComplaint;
DROP TABLE IF EXISTS VendorComplaint;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
