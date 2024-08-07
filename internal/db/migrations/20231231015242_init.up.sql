CREATE TABLE IF NOT EXISTS User (
    id Int NOT NULL AUTO_INCREMENT,
    name Varchar(255) NOT NULL,
    email Varchar(255) NOT NULL UNIQUE,
    emailHash Varchar(64) NOT NULL,
    imageUrl Varchar(255),
    verified TINYINT(1) NOT NULL DEFAULT 0 COMMENT '0: not verified, 1: verified',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX(id, name, emailHash)
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

CREATE TABLE IF NOT EXISTS Admin (
    id Int NOT NULL AUTO_INCREMENT,
    username Varchar(255) NOT NULL UNIQUE,
    password Varchar(255) NOT NULL,
    role Enum('admin', 'staff') NOT NULL DEFAULT 'staff',
    usernameHash Varchar(64) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(id, usernameHash, role)
);

CREATE TABLE IF NOT EXISTS Session (
    id Int NOT NULL AUTO_INCREMENT,
    token Varchar(255) NOT NULL,
    status Enum('online', 'offline') NOT NULL DEFAULT 'online',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    CONSTRAINT unique_token_online UNIQUE (token, (CASE WHEN status = 'online' THEN 1 ELSE NULL END)),
    INDEX(token, status)
);

CREATE TABLE IF NOT EXISTS Blacklist (
    id Int NOT NULL AUTO_INCREMENT,
    token Varchar(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(token)
);

CREATE TABLE IF NOT EXISTS Tag (
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
    description Varchar(255) NOT NULL,
    rate Double NOT NULL,
    latitude Decimal(12, 10) NOT NULL,
    longitude Decimal(13, 10) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES User(id),
    INDEX(id, vendorId)
);

CREATE TABLE IF NOT EXISTS ServiceTag (
    id Int NOT NULL AUTO_INCREMENT,
    serviceId Int NOT NULL,
    tagId Int NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(serviceId) REFERENCES Service(id),
    FOREIGN KEY(tagId) REFERENCES Tag(id),
    INDEX(id, serviceId, tagId)
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

CREATE TABLE IF NOT EXISTS ServicePhoto (
    id Int NOT NULL AUTO_INCREMENT,
    serviceId Int NOT NULL,
    vendorId Int NOT NULL,
    url Varchar(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(serviceId) REFERENCES Service(id) ON DELETE CASCADE,
    FOREIGN KEY(vendorId) REFERENCES User(id) ON DELETE CASCADE,
    INDEX(id, serviceId, vendorId)
);

CREATE TABLE IF NOT EXISTS FrontId (
    id Int NOT NULL AUTO_INCREMENT,
    url Varchar(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(id)
);

CREATE TABLE IF NOT EXISTS BackId (
    id Int NOT NULL AUTO_INCREMENT,
    url Varchar(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(id)
);

CREATE TABLE IF NOT EXISTS Face (
    id Int NOT NULL AUTO_INCREMENT,
    url Varchar(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(id)
);


CREATE TABLE IF NOT EXISTS IdentityVerification (
    id Int NOT NULL AUTO_INCREMENT,
    user Int NOT NULL,
    name Varchar(255) NOT NULL,
    address Varchar(255) NOT NULL,
    idType Varchar(255) NOT NULL,
    idNumber Varchar(255) NOT NULL,
    frontId Int NOT NULL,
    backId Int NOT NULL,
    face Int NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id, user, idType),
    FOREIGN KEY(user) REFERENCES User(id),
    FOREIGN KEY(frontId) REFERENCES FrontId(id),
    FOREIGN KEY(backId) REFERENCES BackId(id),
    FOREIGN KEY(face) REFERENCES Face(id),
    INDEX(id)
);

CREATE TABLE IF NOT EXISTS VendorComplaint (
    id Int NOT NULL AUTO_INCREMENT,
    vendorId Int NOT NULL,
    title Varchar(255) NOT NULL,
    content Text NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES Vendor(id) ON DELETE CASCADE,
    INDEX(id, vendorId)
);

create table if not exists SystemComplaint (
    id Int NOT NULL AUTO_INCREMENT,
    title Varchar(255) NOT NULL,
    detail Text NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    INDEX(id, title)
);

create table if not exists SystemComplaintImage (
    id Int NOT NULL AUTO_INCREMENT,
    complaintId Int NOT NULL,
    url Varchar(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(complaintId) REFERENCES SystemComplaint(id),
    INDEX(id, complaintId)
);

CREATE TABLE IF NOT EXISTS Review (
    id Int NOT NULL AUTO_INCREMENT,
    serviceId Int NOT NULL,
    rating Int NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(serviceId) REFERENCES Service(id) ON DELETE CASCADE,
    INDEX(id, serviceId)
);

CREATE TABLE IF NOT EXISTS Transaction (
    id INT NOT NULL AUTO_INCREMENT,
    vendorId INT NOT NULL,
    clientId INT NOT NULL,
    serviceId INT NOT NULL,
    status Enum('ongoing', 'done', 'cancelled') NOT NULL DEFAULT 'ongoing',
    start TIMESTAMP NOT NULL,
    end TIMESTAMP NOT NULL,
    isReviewed TINYINT(1) NOT NULL DEFAULT 0 COMMENT '0: not reviewed, 1: reviewed',
    isReported TINYINT(1) NOT NULL DEFAULT 0 COMMENT '0: not reported, 1: reported',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(vendorId) REFERENCES Vendor(id) ON DELETE CASCADE,
    FOREIGN KEY(serviceId) REFERENCES Service(id) ON DELETE CASCADE,
    FOREIGN KEY(clientId) REFERENCES User(id)
);

CREATE TABLE IF NOT EXISTS Application (
    id INT NOT NULL AUTO_INCREMENT,
    applicantId INT NOT NULL UNIQUE,
    job Varchar(255) NOT NULL,
    latitude Double NOT NULL,
    longitude Double NOT NULL,
    status Enum('pending', 'rejected', 'approved') NOT NULL DEFAULT 'pending',
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(applicantId) REFERENCES User(id)
);

CREATE TABLE IF NOT EXISTS ApplicationProof (
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

DELIMITER //

CREATE TRIGGER update_vendor_rating
AFTER INSERT ON Review
FOR EACH ROW
BEGIN
    DECLARE avg_rating DECIMAL(5,1);
    
    SELECT ROUND(AVG(rating), 1) INTO avg_rating
    FROM Review
    WHERE serviceId = NEW.serviceId;

    UPDATE Vendor
    SET rating = avg_rating
    WHERE vendorId = (SELECT vendorId FROM Service WHERE id = NEW.serviceId);
END //

DELIMITER ;
