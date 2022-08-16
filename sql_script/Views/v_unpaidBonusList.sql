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
       OR P.payStatus = false
