DROP PROCEDURE IF EXISTS p_GetPayout;

DELIMITER //

CREATE PROCEDURE p_GetPayout(IN _payoutID INT)
BEGIN
    -- Calculate related total campaign budget
    DECLARE _totalBudget DECIMAL(20,2) DEFAULT 0;
    
    SELECT totalBudget INTO _totalBudget
    FROM Payout P
    INNER JOIN (
        SELECT regionCode, payDate, SUM(budget) AS totalBudget
        FROM Campaign
        WHERE approvalStatus IS TRUE
        GROUP BY regionCode, payDate
    ) C ON P.payDate = C.payDate AND P.regionCode = C.regionCode
    WHERE P.id = _payoutID;
    
    -- Calculate related total bonus amount by bonusType
    DROP TEMPORARY TABLE IF EXISTS _bonus;
    
    CREATE TEMPORARY TABLE _bonus
    SELECT SUM(IF(B.bonusType = 0, B.amount, 0)) AS fixedBonus,
           SUM(IF(B.bonusType = 1, B.amount, 0)) AS variableBonus, 
           SUM(IF(B.bonusType = 2, B.amount, 0)) AS transportation,
           SUM(IF(B.bonusType = 3, B.amount, 0)) AS addon,
           SUM(IF(B.bonusType = 4, B.amount, 0)) AS deduction,
           SUM(B.amount) AS total
    FROM PayoutDetail D
    INNER JOIN CampaignBonus B ON D.bonusID = B.id
    WHERE D.payoutID = _payoutID;

    -- Output the result
    SELECT P.id AS payoutID, P.payDate, P.regionCode, B.fixedBonus, B.variableBonus, B.transportation, B.addon, B.deduction, B.total,
           _totalBudget AS totalBudget, _totalBudget - B.total AS difference, P.payStatus AS isPaid, IF(P.payStatus, 2, 1) AS status
    FROM Payout P, _bonus B
    WHERE P.id = _payoutID;
    
    -- Release temp table
    DROP TEMPORARY TABLE _bonus;
END //

DELIMITER ;
