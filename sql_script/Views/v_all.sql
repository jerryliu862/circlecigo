CREATE OR REPLACE VIEW v_campaignList
AS
    SELECT DISTINCT C.id, C.title, C.regionCode, C.startTime, C2.endTime, C2.budget, C2.totalBonus, C2.budget - C2.totalBonus AS bonusDiff, C.remark, C.syncTime, C.approvalStatus, C.approvalTime, C2.approverName
    FROM Campaign C
    INNER JOIN (
        SELECT id, endTime -1 AS endTime, IFNULL(budget, 0) AS budget, fn_GetCampaignTotalBonus(id) AS totalBonus, fn_GetUserName(approverUID) AS approverName
        FROM Campaign
    ) C2 ON C.id = C2.id;


CREATE OR REPLACE VIEW v_regionList
AS
    SELECT R.code, R.name, C.code AS currencyCode, C.name AS currencyName, C.format AS currencyFormat
    FROM Region R
    LEFT JOIN Currency C ON R.currencyCode = C.code
    WHERE R.show IS TRUE;


CREATE OR REPLACE VIEW v_regionTaxList
AS
    SELECT R.code, R.name, C.code AS currencyCode, C.name AS currencyName, C.format AS currencyFormat, T.payType, T.taxFrom, T.taxRate
    FROM Region R
    LEFT JOIN Currency C ON R.currencyCode = C.code
    LEFT JOIN TaxRate T ON R.code = T.regionCode
    WHERE R.show IS TRUE;


CREATE OR REPLACE VIEW v_ungroupedPayout
AS
    SELECT IFNULL(B.payDate, C.payDate) AS payDate, C.regionCode, B.id
    FROM Campaign C
    INNER JOIN CampaignLeaderboard L ON C.id = L.campaignID
    INNER JOIN CampaignRank R ON L.id = R.leaderboardID
    INNER JOIN CampaignBonus B ON R.id = B.rankID
    LEFT JOIN PayoutDetail D ON B.id = D.bonusID
    WHERE C.approvalStatus IS TRUE
      AND D.payoutID IS NULL;


-- v_unpaidBonusList: Streamers and Campaigns list, which is the possible target to add payout adjustment.
CREATE OR REPLACE VIEW v_unpaidBonusList
AS
    SELECT DISTINCT R.id AS rankID, BIN_TO_UUID(S.userID) AS userID, S.openID, C.id AS campaignID, C.title AS campaignTitle, RG.currencyCode AS campaignCurrency, IFNULL(S.regionCode, C.regionCode) AS regionCode
    FROM CampaignRank R
    INNER JOIN CampaignLeaderboard L ON R.leaderboardID = L.id
    INNER JOIN Campaign C ON L.campaignID = C.id
    INNER JOIN Region RG ON C.regionCode = RG.code
    INNER JOIN CampaignBonus B ON R.id = B.rankID
    LEFT JOIN PayoutDetail D ON B.id = D.bonusID
    LEFT JOIN Payout P ON D.payoutID = P.id
    LEFT JOIN Streamer S ON R.userID = S.userID
    WHERE C.approvalStatus = true
      AND P.id IS NULL
       OR P.payStatus = false;


-- v_payoutList: Payout detail list, only include grouped payout data.
CREATE OR REPLACE VIEW v_payoutDetail
AS
    SELECT P.id AS payoutID, R.id AS rankID, S.openID, BIN_TO_UUID(S.userID) AS userID, C.regionCode, S.regionCode AS streamerRegion, C.id AS campaignID, C.title AS campaignTitle, B.id AS bonusID, B.bonusType, B.amount, B.remark
    FROM Payout P
    INNER JOIN PayoutDetail D ON P.id = D.payoutID
    INNER JOIN CampaignBonus B ON D.bonusID = B.id
    INNER JOIN CampaignRank R ON B.rankID = R.id
    INNER JOIN CampaignLeaderboard L ON R.leaderboardID = L.id
    INNER JOIN Campaign C ON L.campaignID = C.id
    INNER JOIN Streamer S ON R.userID = S.userID
    ORDER BY amount = 0, payoutID, rankID;


-- v_payoutList: Payout detail list, only include not yet grouped payout data.
CREATE OR REPLACE VIEW v_payoutDetail2
AS
    SELECT IFNULL(B.payDate, C.payDate) AS payDate, S.openID, BIN_TO_UUID(S.userID) AS userID, C.regionCode, S.regionCode AS streamerRegion, C.id AS campaignID, C.title AS campaignTitle, B.id AS bonusID, B.bonusType, B.amount, B.remark
    FROM Campaign C
    INNER JOIN CampaignLeaderboard L ON C.id = L.campaignID
    INNER JOIN CampaignRank R ON L.id = R.leaderboardID
    INNER JOIN Streamer S ON R.userID = S.userID
    INNER JOIN CampaignBonus B ON R.id = B.rankID
    LEFT JOIN PayoutDetail D ON B.id = D.bonusID
    WHERE C.approvalStatus IS TRUE
      AND D.payoutID IS NULL
    ORDER BY amount = 0, payDate, regionCode, streamerRegion, campaignID, bonusType;


CREATE OR REPLACE VIEW v_payoutRecord
 AS
    SELECT payoutID, userID, agencyID, accountID, payType, payAmount, taxableAmount, taxRate, IF(taxFrom IS NULL, NULL, IF(taxableAmount >= taxFrom, taxableAmount * taxRate, 0)) AS taxAmount, accountName, taxID, feeOwnerType
    FROM (
        SELECT P.id AS payoutID, S.userID, S.agencyID, S.accountID, K.payType, 
            SUM(B.amount) AS payAmount, SUM(IF(B.bonusType != 2, B.amount, 0)) AS taxableAmount, T.taxFrom, T.taxRate, K.accountName, K.payeeID AS taxID, K.feeOwnerType
        FROM Payout P
        INNER JOIN PayoutDetail D ON P.id = D.payoutID
        INNER JOIN CampaignBonus B ON D.bonusID = B.id
        INNER JOIN CampaignRank R ON B.rankID = R.id
        INNER JOIN (
            SELECT S.userID, S.agencyID, S.regionCode, IF(A.ID IS NOT NULL, A.accountID, S.accountID) AS accountID
            FROM Streamer S
            LEFT JOIN StreamerAgency A ON S.agencyID = A.id
        ) S ON R.userID = S.userID
        LEFT JOIN BankAccount K ON S.accountID = K.id
        LEFT JOIN TaxRate T ON S.regionCode = T.regionCode AND K.payType = T.payType
        GROUP BY P.id, S.userID, S.agencyID, S.accountID, K.payType, T.taxFrom, T.taxRate, K.payeeID, K.feeOwnerType
    ) A;


-- v_payoutReportList: Payout report data list, only include grouped payout data.
CREATE OR REPLACE VIEW v_payoutReportList
AS
    SELECT P.id AS payoutID, P.payDate, P.regionCode, 
        SUM(IF(R.payType = 'F', R.payAmount, 0)) AS foreignAmount,
        SUM(IF(R.payType = 'L', R.payAmount, 0)) AS localAmount,
        SUM(IF(R.payType = 'A', R.payAmount, 0)) AS agencyAmount,
        COUNT(DISTINCT R.userID) AS streamerCount,
        COUNT(DISTINCT CASE WHEN R.accountID IS NULL THEN R.userID END) missingCount,
        SUM(R.payAmount) AS totalAmount,
        SUM(IFNULL(R.taxAmount, 0)) AS totalTaxAmount,
        P.payStatus
    FROM Payout P
    INNER JOIN v_payoutRecord R ON P.id = R.payoutID
    GROUP BY P.id, P.payDate, P.regionCode, P.payStatus;


-- v_payoutReportDetail: Payout report data, only include grouped payout data.
CREATE OR REPLACE VIEW v_payoutReportDetail
AS
    SELECT DISTINCT P.id AS payoutID, P.regionCode, BIN_TO_UUID(S.userID) AS userID, S.openID, S.regionCode AS streamerRegion, R.agencyID, R.payType, R.payAmount, R.taxAmount, 
           A.swiftCode, A.bankCode, A.bankName, A.branchCode, A.branchName, A.accountNo, A.accountName, A.feeOwnerType, A.payeeID AS taxID
    FROM Payout P
    INNER JOIN v_payoutRecord R ON P.id = R.payoutID
    INNER JOIN Streamer S ON R.userID = S.userID
    LEFT JOIN BankAccount A ON R.accountID = A.id;


-- v_payoutReportDetailByAgency: Payout report data, only include grouped payout data grouped by agency.
CREATE OR REPLACE VIEW v_payoutReportDetailByAgency
AS
    SELECT P.id AS payoutID, P.regionCode, A.id AS agencyID, IF(A.name IS NULL or A.name = '', R.accountName, A.name) AS agencyName, SUM(R.payAmount) AS totalAmount, SUM(R.taxAmount) AS totalTaxAmount, 
        K.bankCode, K.bankName, K.branchCode, K.branchName, K.accountNo, K.accountName, K.feeOwnerType, R.taxID
    FROM Payout P
    INNER JOIN v_payoutRecord R ON P.id = R.payoutID
    INNER JOIN Streamer S ON R.userID = S.userID
    INNER JOIN StreamerAgency A ON S.agencyID = A.id
    LEFT JOIN BankAccount K ON A.accountID = K.id
    WHERE R.payType = 'A'
    GROUP BY P.id, P.regionCode, A.id, agencyName, K.bankCode, K.bankName, K.branchCode, K.branchName, K.accountNo, K.accountName, K.feeOwnerType, R.taxID;


-- v_payoutReportDetailByAgency2: Payout report data, only include grouped payout data grouped by agency.
CREATE OR REPLACE VIEW v_payoutReportDetailByAgency2
AS
    SELECT P.regionCode, O.payoutID, O.agencyID, BIN_TO_UUID(S.userID) AS userID, S.openID, S.regionCode AS streamerRegion, C.title AS campaignTitle, O.accountName, SUM(B.Amount) AS payAmount
    FROM Payout P
    INNER JOIN v_payoutRecord O ON P.id = O.payoutID
    INNER JOIN PayoutDetail D ON P.id = D.payoutID
    INNER JOIN CampaignBonus B ON D.bonusID = B.id
    INNER JOIN CampaignRank R ON B.rankID = R.id AND O.userID = R.userID
    INNER JOIN CampaignLeaderboard L ON R.leaderboardID = L.id
    INNER JOIN Campaign C ON L.campaignID = C.id
    INNER JOIN Streamer S ON R.userID = S.userID
    WHERE O.payType = 'A'
    GROUP BY payoutID, agencyID, userID, openID, streamerRegion, accountName, campaignTitle;


CREATE OR REPLACE VIEW v_payoutReportDetailByBonusType
AS
    SELECT P.id AS payoutID, P.regionCode, BIN_TO_UUID(S.userID) AS userID, S.openID, S.regionCode AS streamerRegion, B.bonusType, SUM(B.amount) AS payAmount, 0 AS taxAmount, K.payType, 
           K.bankCode, K.bankName, K.branchCode, K.branchName, K.accountNo, K.accountName, K.feeOwnerType, K.payeeID AS taxID
    FROM Payout P
    INNER JOIN PayoutDetail D ON P.id = D.payoutID
    INNER JOIN CampaignBonus B ON D.bonusID = B.id
    INNER JOIN CampaignRank R ON B.rankID = R.id
    INNER JOIN (
        SELECT S.userID, S.openID, S.regionCode, IF(A.ID IS NOT NULL, A.accountID, S.accountID) AS accountID
        FROM Streamer S
        LEFT JOIN StreamerAgency A ON S.agencyID = A.id
    ) S ON R.userID = S.userID
    LEFT JOIN BankAccount K ON S.accountID = K.id
    GROUP BY P.id, P.regionCode, S.userID, S.openID, S.regionCode, B.bonusType, K.payType, K.bankCode, K.bankName, K.branchCode, K.branchName, K.accountNo, K.accountName, K.feeOwnerType, K.payeeID;
