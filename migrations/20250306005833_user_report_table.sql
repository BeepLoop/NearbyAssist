-- +goose Up
CREATE TABLE IF NOT EXISTS UserReport (
    id VARCHAR(255) NOT NULL,
    reporterUserId VARCHAR(255) NOT NULL,
    reportedUserId VARCHAR(255) NOT NULL,
    category ENUM('misconduct', 'booking_related') NOT NULL,
    bookingId VARCHAR(255),
    reason VARCHAR(255) NOT NULL,
    detail TEXT,
    status ENUM('pending', 'resolved', 'dismissed') DEFAULT 'pending',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(reporterUserId) REFERENCES User(id) ON DELETE CASCADE,
    FOREIGN KEY(reportedUserId) REFERENCES User(id) ON DELETE CASCADE,
    FOREIGN KEY(bookingId) REFERENCES Booking(id),
    INDEX(id, reportedUserId)
);

CREATE TABLE IF NOT EXISTS UserReportImage (
    id VARCHAR(255) NOT NULL,
    reportId VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(reportId) REFERENCES UserReport(id),
    INDEX(id, reportId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS UserReportImage;
DROP TABLE IF EXISTS UserReport;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
