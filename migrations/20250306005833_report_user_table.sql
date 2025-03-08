-- +goose Up
CREATE TABLE IF NOT EXISTS ReportedUser (
    id VARCHAR(255) NOT NULL,
    userId VARCHAR(255) NOT NULL,
    reason VARCHAR(255) NOT NULL,
    detail TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(userId) REFERENCES User(id) ON DELETE CASCADE,
    INDEX(id, userId)
);

CREATE TABLE IF NOT EXISTS ReportedUserImage (
    id VARCHAR(255) NOT NULL,
    reportId VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(reportId) REFERENCES ReportedUser(id),
    INDEX(id, reportId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS ReportedUserImage;
DROP TABLE IF EXISTS ReportedUser;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
