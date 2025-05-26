-- +goose Up
CREATE TABLE IF NOT EXISTS VendorExpertise(
    vendorId VARCHAR(255) NOT NULL,
    expertiseId VARCHAR(255) NOT NULL,
    applicationId VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    supportingImage VARCHAR(255) NOT NULL,
    FOREIGN KEY(vendorId) REFERENCES Vendor(vendorId),
    FOREIGN KEY(expertiseId) REFERENCES Expertise(id),
    FOREIGN KEY(applicationId) REFERENCES Application(id),
    FOREIGN KEY(supportingImage) REFERENCES SupportingImage(id),
    INDEX(vendorId, expertiseId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS VendorExpertise;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
