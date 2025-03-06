-- +goose Up
CREATE TABLE IF NOT EXISTS BugReport (
    id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    detail Text NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(id, title)
);

CREATE TABLE IF NOT EXISTS BugReportImage (
    id VARCHAR(255) NOT NULL,
    reportId VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(reportId) REFERENCES BugReport(id),
    INDEX(id, reportId)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP DATABASE IF EXISTS BugReportImage;
DROP DATABASE IF EXISTS BugReport;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
