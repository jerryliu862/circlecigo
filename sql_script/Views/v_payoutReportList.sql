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
    GROUP BY P.id, P.payDate, P.regionCode, P.payStatus
