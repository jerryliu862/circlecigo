DROP PROCEDURE IF EXISTS p_GetCampaignDetailByRank;

DELIMITER //

CREATE PROCEDURE p_GetCampaignDetailByRank(IN leaderboardID VARBINARY(16))
BEGIN
    SELECT R.leaderboardID, R.id AS rankID, R.rank, R.score, BIN_TO_UUID(S.userID) AS userID, S.openID, S.regionCode AS streamerRegion, 
           SUM(IF(B.bonusType = 0, B.amount, 0)) AS fixedBonus, SUM(IF(B.bonusType = 1, B.amount, 0)) AS variableBonus, SUM(IF(B.bonusType IN (0, 1), B.amount, 0)) AS totalBonus
    FROM CampaignRank R
    LEFT JOIN CampaignBonus B ON R.id = B.rankID
    LEFT JOIN Streamer S ON R.userID = S.userID
    WHERE R.leaderboardID = leaderboardID
    GROUP BY R.leaderboardID, R.id, R.rank, R.score, S.userID, S.regionCode
    ORDER BY R.leaderboardID, R.rank;
END //

DELIMITER ;
