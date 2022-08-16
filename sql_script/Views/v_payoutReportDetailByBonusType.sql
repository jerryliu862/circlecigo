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
    GROUP BY P.id, P.regionCode, S.userID, S.openID, S.regionCode, B.bonusType, K.payType, K.bankCode, K.bankName, K.branchCode, K.branchName, K.accountNo, K.accountName, K.feeOwnerType, K.payeeID
