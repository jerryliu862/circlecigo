-- v_payoutReportDetail: Payout report data, only include grouped payout data.
CREATE OR REPLACE VIEW v_payoutReportDetail
AS
    SELECT DISTINCT P.id AS payoutID, P.regionCode, BIN_TO_UUID(S.userID) AS userID, S.openID, S.regionCode AS streamerRegion, R.agencyID, R.payType, R.payAmount, R.taxAmount, 
           A.swiftCode, A.bankCode, A.bankName, A.branchCode, A.branchName, A.accountNo, A.accountName, A.feeOwnerType, A.payeeID AS taxID
    FROM Payout P
    INNER JOIN v_payoutRecord R ON P.id = R.payoutID
    INNER JOIN Streamer S ON R.userID = S.userID
    LEFT JOIN BankAccount A ON R.accountID = A.id
