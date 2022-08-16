DROP PROCEDURE IF EXISTS p_GetCampaignDetailByStage;

DELIMITER //

CREATE PROCEDURE p_GetCampaignDetailByStage(IN campaignID INT)
BEGIN
    SELECT L.id AS leaderboardID, L.title, SUM(IF(B.bonusType = 0, B.amount, 0)) AS totalFixedBonus, SUM(IF(B.bonusType = 1, B.amount, 0)) AS totalVariableBonus, SUM(IF(B.bonusType IN (0, 1), B.amount, 0)) AS totalBonus
    FROM CampaignLeaderboard L
    INNER JOIN CampaignRank R ON L.id = R.leaderboardID
    LEFT JOIN CampaignBonus B ON R.id = B.rankID
    WHERE L.campaignID = campaignID
    GROUP BY L.campaignID, L.id, L.title
    ORDER BY L.title;
END //

DELIMITER ;
