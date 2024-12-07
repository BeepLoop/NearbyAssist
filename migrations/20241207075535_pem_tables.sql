-- +goose Up
CREATE TABLE IF NOT EXISTS PublicKey (
    id VARCHAR(255) NOT NULL,
    owner VARCHAR(255) NOT NULL,
    pem TEXT NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(owner) REFERENCES User(id),
    INDEX(id, owner)
);

CREATE TABLE IF NOT EXISTS PrivateKey (
    id VARCHAR(255) NOT NULL,
    owner VARCHAR(255) NOT NULL,
    pem TEXT NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(owner) REFERENCES User(id),
    INDEX(id, owner)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS PublicKey;
DROP TABLE IF EXISTS PrivateKey;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
