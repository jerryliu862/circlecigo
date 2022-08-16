/* DROP All tables (except User, UserAuth) if necessary */
/*
 DROP TABLE IF EXISTS PayoutRecord;
 DROP TABLE IF EXISTS PayoutDetail;
 DROP TABLE IF EXISTS Payout;
 DROP TABLE IF EXISTS CampaignBonus;
 DROP TABLE IF EXISTS CampaignRank;
 DROP TABLE IF EXISTS CampaignLeaderboard;
 DROP TABLE IF EXISTS Campaign;
 DROP TABLE IF EXISTS Streamer;
 DROP TABLE IF EXISTS StreamerAgency;
 DROP TABLE IF EXISTS BankAccount;

 DROP TABLE IF EXISTS TaxRate;
 DROP TABLE IF EXISTS Region;
 DROP TABLE IF EXISTS Currency;
*/

/* Currency and Region */

CREATE TABLE IF NOT EXISTS Currency (
    `code`   VARCHAR(3)  NOT NULL,
    `name`   VARCHAR(50) NOT NULL,
    `format` TINYINT     NOT NULL,
    PRIMARY KEY (code)
);

CREATE TABLE IF NOT EXISTS Region (
    `code`         VARCHAR(10) NOT NULL,
    `name`         VARCHAR(50) NOT NULL,
    `currencyCode` VARCHAR(3),
    `show`         BOOLEAN     NOT NULL,
    PRIMARY KEY (code),
    FOREIGN KEY (currencyCode) REFERENCES Currency (code)
);

/* User */

CREATE TABLE IF NOT EXISTS User (
    `id`          INT          AUTO_INCREMENT,
    `googleID`    VARCHAR(255) NOT NULL,
    `name`        VARCHAR(255) NOT NULL,
    `email`       VARCHAR(50)  NOT NULL,
    `status`      TINYINT      NOT NULL,
    `createTime`  DATETIME     NOT NULL,
    `creatorUID`  INT          NOT NULL,
    `modifyTime`  DATETIME,
    `modifierUID` INT,
    PRIMARY KEY (id),
    UNIQUE (email),
    FOREIGN KEY (creatorUID)  REFERENCES `User` (id),
    FOREIGN KEY (modifierUID) REFERENCES `User` (id)
);

CREATE TABLE IF NOT EXISTS UserAuth (
    `id`         INT         AUTO_INCREMENT,
    `uid`        INT         NOT NULL,
    `regionCode` VARCHAR(10) NOT NULL,
    `authType`   TINYINT     NOT NULL,
    `authLevel`  TINYINT     NOT NULL,
    `sendNotify` BOOLEAN     NOT NULL DEFAULT TRUE,
    `createTime` DATETIME    NOT NULL,
    `creatorUID` INT         NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (uid)        REFERENCES User (id),
    FOREIGN KEY (regionCode) REFERENCES Region (code),
    FOREIGN KEY (creatorUID) REFERENCES User (id)
);

CREATE TABLE IF NOT EXISTS TaxRate (
    `id`         INT          AUTO_INCREMENT,
    `regionCode` VARCHAR(10)  NOT NULL,
    `payType`    CHAR(1)      NOT NULL,
    `taxFrom`    INT          NOT NULL,
    `taxRate`    DECIMAL(3,3) NOT NULL,
    `createTime`  DATETIME    NOT NULL,
    `creatorUID`  INT         NOT NULL,
    `modifyTime`  DATETIME,
    `modifierUID` INT,
    PRIMARY KEY (id),
    FOREIGN KEY (regionCode)  REFERENCES Region (code),
    FOREIGN KEY (creatorUID)  REFERENCES User (id),    
    FOREIGN KEY (modifierUID) REFERENCES User (id),
    CONSTRAINT u_region_and_payType UNIQUE (regionCode, payType)
);

/* Bank */

CREATE TABLE IF NOT EXISTS BankAccount (
    `id`           INT          AUTO_INCREMENT,
    `swiftCode`    VARCHAR(20),
    `bankCode`     VARCHAR(50),
    `bankName`     VARCHAR(200),
    `branchCode`   VARCHAR(50),
    `branchName`   VARCHAR(200),
    `accountNo`    VARCHAR(50)  NOT NULL,
    `accountName`  VARCHAR(100) NOT NULL,
    `accountType`  TINYINT      NOT NULL,
    `payeeType`    TINYINT      NOT NULL,
    `payeeID`      VARCHAR(10),
    `ownerType`    TINYINT,
    `feeOwnerType` TINYINT,
    `payType`      CHAR(1)      GENERATED ALWAYS AS (CASE WHEN ownerType = 1 AND accountType = 1 THEN 'L' WHEN ownerType = 1 AND accountType = 2 THEN 'F' WHEN ownerType = 2 AND accountType = 1 THEN 'A' ELSE NULL END),
    `syncTime`     DATETIME     NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT u_bank_and_account UNIQUE (swiftCode, bankCode, branchCode, accountNo)
);

/* Streamer */

CREATE TABLE IF NOT EXISTS StreamerAgency (
    `id`        INT           NOT NULL,
    `name`      VARCHAR(50),
    `accountID` INT,
    `syncTime`  DATETIME      NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (accountID) REFERENCES BankAccount (id)
);

CREATE TABLE IF NOT EXISTS Streamer (
    `userID`     BINARY(16)   NOT NULL,
    `openID`     VARCHAR(50),
    `agencyID`   INT,
    `name`       VARCHAR(50),
    `regionCode` VARCHAR(10),
    `accountID`  INT,
    `syncTime`   DATETIME     NOT NULL,
    PRIMARY KEY (userID),
    FOREIGN KEY (agencyID)   REFERENCES StreamerAgency (id),
    FOREIGN KEY (regionCode) REFERENCES Region (code),
    FOREIGN KEY (accountID)  REFERENCES BankAccount (id)
);

/* Campaign */

CREATE TABLE IF NOT EXISTS Campaign (
    `id`             INT,
    `title`          VARCHAR(255)  NOT NULL,
    `regionCode`     VARCHAR(10),
    `regionList`     VARCHAR(500),
    `startTime`      BIGINT        NOT NULL,
    `endTime`        BIGINT        NOT NULL,
    `payDate`        DATE          GENERATED ALWAYS AS (DATE(DATE_FORMAT(FROM_UNIXTIME(endTime), '%Y-%m-01'))),
    `budget`         DECIMAL(20,0),
    `totalBonus`     DECIMAL(19,2),
    `remark`         VARCHAR(100),
    `syncTime`       DATETIME      NOT NULL,
    `modifyTime`     DATETIME,
    `modifierUID`    INT,
    `approvalStatus` BOOLEAN       NOT NULL,
    `approvalTime`   DATETIME,
    `approverUID`    INT,
    PRIMARY KEY (id),
    FOREIGN KEY (regionCode)  REFERENCES Region (code),
    FOREIGN KEY (modifierUID) REFERENCES User (id),
    FOREIGN KEY (approverUID) REFERENCES User (id)
);

CREATE TABLE IF NOT EXISTS CampaignLeaderboard (
    `id`         BINARY(16)   NOT NULL,
    `campaignID` INT          NOT NULL,
    `title`      VARCHAR(255) NOT NULL,
    `syncTime`   DATETIME     NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (campaignID) REFERENCES Campaign (id)
);

CREATE TABLE IF NOT EXISTS CampaignRank (
    `id`            INT       AUTO_INCREMENT,
    `leaderboardID` BINARY(16) NOT NULL,
    `rank`          TINYINT    NOT NULL,
    `score`         BIGINT     NOT NULL,
    `userID`        BINARY(16) NOT NULL,
    `syncTime`      DATETIME   NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (leaderboardID) REFERENCES CampaignLeaderboard (id),
    FOREIGN KEY (userID) REFERENCES Streamer (userID),
    CONSTRAINT u_leaderboard_and_streamer UNIQUE (leaderboardID, userID)
);

CREATE TABLE IF NOT EXISTS CampaignBonus (
    `id`          INT           AUTO_INCREMENT,
    `rankID`      INT           NOT NULL,
    `bonusType`   TINYINT       NOT NULL,
    `amount`      DECIMAL(20,2) NOT NULL,
    `payDate`     DATE,
    `remark`      VARCHAR(100),
    `createTime`  DATETIME      NOT NULL,
    `creatorUID`  INT           NOT NULL,
    `modifyTime`  DATETIME,
    `modifierUID` INT,
    PRIMARY KEY (id),
    FOREIGN KEY (rankID)      REFERENCES CampaignRank (id),
    FOREIGN KEY (creatorUID)  REFERENCES User (id),
    FOREIGN KEY (modifierUID) REFERENCES User (id)
);

/* Payout */

CREATE TABLE IF NOT EXISTS Payout (
    `id`          INT         AUTO_INCREMENT,
    `regionCode`  VARCHAR(10) NOT NULL,
    `payDate`     DATE        NOT NULL,
    `payStatus`   BOOLEAN     NOT NULL,
    `createTime`  DATETIME    NOT NULL,
    `creatorUID`  INT         NOT NULL,
    `modifyTime`  DATETIME,
    `modifierUID` INT,
    PRIMARY KEY (id),
    FOREIGN KEY (regionCode)  REFERENCES Region (code),
    FOREIGN KEY (creatorUID)  REFERENCES User (id),
    FOREIGN KEY (modifierUID) REFERENCES User (id)
);

CREATE TABLE IF NOT EXISTS PayoutDetail (
    `payoutID`   INT      NOT NULL,
    `bonusID`    INT      NOT NULL,
    `createTime` DATETIME NOT NULL,
    `creatorUID` INT      NOT NULL,
    PRIMARY KEY (payoutID, bonusID),
    FOREIGN KEY (payoutID)   REFERENCES Payout (id),
    FOREIGN KEY (bonusID)    REFERENCES CampaignBonus (id),
    FOREIGN KEY (creatorUID) REFERENCES User (id)
);

CREATE TABLE IF NOT EXISTS PayoutRecord (
    `payoutID`      INT           NOT NULL,
    `userID`        BINARY(16)    NOT NULL,
    `agencyID`      INT,
    `accountID`     INT,
    `payType`       CHAR(1),
    `payAmount`     DECIMAL(20,2) NOT NULL,
    `taxableAmount` DECIMAL(20,2) NOT NULL,
    `taxRate`       DECIMAL(3,3),
    `taxAmount`     DECIMAL(18,2),
    `taxID`         VARCHAR(10),
    `feeOwnerType`  TINYINT,
    `createTime`    DATETIME      NOT NULL,
    `creatorUID`    INT           NOT NULL,
    PRIMARY KEY (payoutID, userID),
    FOREIGN KEY (payoutID)    REFERENCES Payout (id),
    FOREIGN KEY (userID)      REFERENCES Streamer (userID),
    FOREIGN KEY (accountID)   REFERENCES BankAccount (id),
    FOREIGN KEY (agencyID)    REFERENCES StreamerAgency (id),
    FOREIGN KEY (creatorUID)  REFERENCES User (id)
);
