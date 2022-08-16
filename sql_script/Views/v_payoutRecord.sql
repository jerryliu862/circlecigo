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
    ) A
