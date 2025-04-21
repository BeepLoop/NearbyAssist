-- +goose Up
CREATE TABLE IF NOT EXISTS UserIdentification (
    userId VARCHAR(255) NOT NULL,
    identificationId VARCHAR(255) NOT NULL,
    FOREIGN KEY(userId) REFERENCES User(id),
    FOREIGN KEY(identificationId) REFERENCES Identification(id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS UserIdentification;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
