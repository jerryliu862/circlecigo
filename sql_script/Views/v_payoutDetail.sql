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
    ORDER BY amount = 0, payoutID, rankID
