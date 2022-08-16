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
    ORDER BY amount = 0, payDate, regionCode, streamerRegion, campaignID, bonusType
