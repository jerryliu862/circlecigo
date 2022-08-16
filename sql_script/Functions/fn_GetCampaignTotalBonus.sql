DROP FUNCTION IF EXISTS fn_GetCampaignTotalBonus;

DELIMITER //
CREATE FUNCTION fn_GetCampaignTotalBonus (
    campaignID INT
) RETURNS DECIMAL(20, 2) DETERMINISTIC
BEGIN
  DECLARE result DECIMAL(20, 2);
  SET result = 0;
  
  SELECT SUM(B.amount)
  INTO result
  FROM Campaign C
  INNER JOIN CampaignLeaderboard L ON C.id = L.campaignID
  INNER JOIN CampaignRank R ON L.id = R.leaderboardID
  INNER JOIN CampaignBonus B ON R.id = B.rankID AND (B.payDate IS NULL OR B.payDate = C.payDate)
  WHERE B.bonusType IN (0, 1)
    AND C.id = campaignID;
    
  RETURN result;
END //

DELIMITER ;
