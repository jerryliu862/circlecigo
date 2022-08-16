CREATE OR REPLACE VIEW v_ungroupedPayout
AS
    SELECT IFNULL(B.payDate, C.payDate) AS payDate, C.regionCode, B.id
    FROM Campaign C
    INNER JOIN CampaignLeaderboard L ON C.id = L.campaignID
    INNER JOIN CampaignRank R ON L.id = R.leaderboardID
    INNER JOIN CampaignBonus B ON R.id = B.rankID
    LEFT JOIN PayoutDetail D ON B.id = D.bonusID
    WHERE C.approvalStatus IS TRUE
      AND D.payoutID IS NULL
