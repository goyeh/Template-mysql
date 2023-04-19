USE news;

CREATE TABLE IF NOT EXISTS `news`.`articles` (
    `uid`         INTEGER NOT NULL AUTO_INCREMENT ,
    `title`       VARCHAR(128) NOT NULL UNIQUE,
    `content`     TEXT NOT NULL ,
    `author`      VARCHAR(80) NOT NULL,
    `email`       VARCHAR(80) NOT NULL,
    `topic`       VARCHAR(512) NOT NULL,
    `cat`         VARCHAR(512) NOT NULL,
    `link`        VARCHAR(256) NOT NULL,
    `detail`      VARCHAR(1024) NOT NULL,
    `rating`      INT NOT NULL DEFAULT 0,
    `created`     TIMESTAMP ,
    PRIMARY KEY (`uid`),
    INDEX (`created`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE IF NOT EXISTS `news`.`genrated` (
    `uid`         INTEGER NOT NULL AUTO_INCREMENT ,
    `title`       VARCHAR(128) NOT NULL UNIQUE,
    `content`     TEXT NOT NULL ,
    `author`      VARCHAR(80) NOT NULL,
    `email`       VARCHAR(80) NOT NULL,
    `topic`       VARCHAR(512) NOT NULL,
    `cat`         VARCHAR(512) NOT NULL,
    `link`        VARCHAR(256) NOT NULL,
    `detail`      VARCHAR(256) NOT NULL,
    `rating`      INT NOT NULL DEFAULT 0,
    `created`     TIMESTAMP ,
    PRIMARY KEY (`uid`),
    INDEX (`created`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;


CREATE TABLE IF NOT EXISTS `news`.`control` (
    `target`      VARCHAR(40) PRIMARY KEY NOT NULL,
    `note`        VARCHAR(128) NOT NULL,
    `live`        TINYINT NOT NULL DEFAULT 1,
    `platform`    VARCHAR(45) NOT NULL,
    `timestamp`   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `lastupdate`  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX (`target`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;


CREATE TABLE IF NOT EXISTS `news`.`accounts` (
    `urn`         INT AUTO_INCREMENT PRIMARY KEY,
    `apikey`      VARCHAR(64) NOT NULL UNIQUE,
    `email`       VARCHAR(228) NOT NULL UNIQUE,
    `otpkey`      VARCHAR(128) DEFAULT '',
    `plan`        VARCHAR(45) DEFAULT 'FREE' ,
    `allocated`   INT DEFAULT 500,
    `used`        INT DEFAULT 0,
    `end`         TIMESTAMP DEFAULT (CURDATE() + INTERVAL 1 MONTH),
    `timestamp`   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `lastupdate`  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX (`apikey`),
    INDEX (`email`)   
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;


DELIMITER //

CREATE TRIGGER trigger_name BEFORE INSERT ON accounts
FOR EACH ROW
BEGIN
  SET NEW.apikey = LEFT(MD5(CONCAT('salt', NEW.email)), 24);
END//

DELIMITER ;
