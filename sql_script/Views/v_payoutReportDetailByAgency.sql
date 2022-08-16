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
    GROUP BY P.id, P.regionCode, A.id, agencyName, K.bankCode, K.bankName, K.branchCode, K.branchName, K.accountNo, K.accountName, K.feeOwnerType, R.taxID
