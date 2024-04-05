CREATE TABLE IF NOT EXISTS User (
    id Int NOT NULL AUTO_INCREMENT,
    name Varchar(255) NOT NULL,
    email Varchar(255) NOT NULL UNIQUE,
    imageUrl Varchar(255),
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX(id, name, email)
);

CREATE TABLE IF NOT EXISTS Session (
    id Int NOT NULL AUTO_INCREMENT,
    userId Int NOT NULL,
    token Text NOT NULL,
    status Enum('online', 'offline') NOT NULL DEFAULT 'online',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(userId) REFERENCES User(id),
    CONSTRAINT unique_userId_online UNIQUE (userId, (CASE WHEN status = 'online' THEN 1 ELSE NULL END)),
    INDEX(userId, status)
);

CREATE TABLE IF NOT EXISTS Category (
    id Int NOT NULL AUTO_INCREMENT,
    title Varchar(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(id, title)
);

CREATE TABLE IF NOT EXISTS Service (
    id Int NOT NULL AUTO_INCREMENT,
    vendorId Int NOT NULL,
    title Varchar(255) NOT NULL,
    description Varchar(255) NOT NULL,
    rate Double NOT NULL,
    location Geometry NOT NULL SRID 4326,
    category  Int NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES User(id),
    FOREIGN KEY(category) REFERENCES Category(id),
    INDEX(id, vendorId, category)
);

CREATE TABLE IF NOT EXISTS Message (
    id Int NOT NULL AUTO_INCREMENT,
    sender Int NOT NULL,
    receiver Int NOT NULL,
    content Text NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(sender) REFERENCES User(id),
    FOREIGN KEY(receiver) REFERENCES User(id),
    INDEX(id, sender, receiver)
);

CREATE TABLE IF NOT EXISTS Vendor (
    id Int NOT NULL AUTO_INCREMENT,
    vendorId Int NOT NULL,
    rating Decimal(5,1) NOT NULL DEFAULT 0.0,
    job Varchar(255) NOT NULL,
    restricted TINYINT(1) NOT NULL DEFAULT 0,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES User(id),
    INDEX(id, vendorId, job)
);

CREATE TABLE IF NOT EXISTS ServicePhoto (
    id Int NOT NULL AUTO_INCREMENT,
    serviceId Int NOT NULL,
    vendorId Int NOT NULL,
    url Varchar(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(serviceId) REFERENCES Service(id),
    FOREIGN KEY(vendorId) REFERENCES Vendor(id),
    INDEX(id, serviceId, vendorId)
);

CREATE TABLE IF NOT EXISTS Complaint (
    id Int NOT NULL AUTO_INCREMENT,
    vendorId Int NOT NULL,
    code Int NOT NULL,
    title Varchar(255) NOT NULL,
    content Text NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES Vendor(id),
    INDEX(id, vendorId)
);

CREATE TABLE IF NOT EXISTS Review (
    id Int NOT NULL AUTO_INCREMENT,
    serviceId Int NOT NULL,
    rating Int NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(serviceId) REFERENCES Service(id),
    INDEX(id, serviceId)
);

CREATE TABLE IF NOT EXISTS Transaction (
    id INT NOT NULL AUTO_INCREMENT,
    vendorId INT NOT NULL,
    clientId INT NOT NULL,
    serviceId INT NOT NULL,
    status Enum('ongoing', 'done', 'cancelled') NOT NULL DEFAULT 'ongoing',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES Vendor(id),
    FOREIGN KEY(serviceId) REFERENCES Service(id),
    FOREIGN KEY(clientId) REFERENCES User(id)
);

create table if not exists Application (
    id INT NOT NULL AUTO_INCREMENT,
    applicantId INT NOT NULL,
    job Varchar(255) NOT NULL,
    latitude Double NOT NULL,
    longitude Double NOT NULL,
    status Enum('pending', 'rejected', 'approved') NOT NULL DEFAULT 'pending',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(applicantId) REFERENCES User(id)
);

create table if not exists ApplicationProof (
    id INT NOT NULL AUTO_INCREMENT,
    applicationId INT NOT NULL,
    applicantId INT NOT NULL,
    url Varchar(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(applicationId) REFERENCES Application(id),
    FOREIGN KEY(applicantId) REFERENCES Application(applicantId)
);

CREATE TRIGGER update_vendor_rating
AFTER INSERT ON Review
FOR EACH ROW
BEGIN
    DECLARE avg_rating DECIMAL(5,1);
    
    -- Compute the average rating for the given serviceId
    SELECT ROUND(AVG(rating), 1) INTO avg_rating
    FROM Review
    WHERE serviceId = NEW.serviceId;

    -- Update the rating field in the Vendor table
    UPDATE Vendor
    SET rating = avg_rating
    WHERE vendorId = (SELECT vendorId FROM Service WHERE id = NEW.serviceId);
END;
