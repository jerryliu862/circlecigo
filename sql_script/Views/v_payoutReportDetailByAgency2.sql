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
    GROUP BY payoutID, agencyID, userID, openID, streamerRegion, accountName, campaignTitle
