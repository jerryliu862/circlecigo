DROP PROCEDURE IF EXISTS p_GetPayoutList;

DELIMITER //

CREATE PROCEDURE p_GetPayoutList(IN _payDate DATE, IN _regionList TEXT, IN _isGrouped BOOL, IN _isPaid BOOL, IN _limitCount INT, IN _offsetCount INT)
BEGIN
    -- Set data counter
    DECLARE _ungroupedCount INT unsigned DEFAULT 0;
    DECLARE _groupedCount INT unsigned DEFAULT 0;
    DECLARE _pagedCount INT unsigned DEFAULT 0;

    -- Modify variable
    SET _regionList = REPLACE(_regionList, '|', ',');
    
    -- Construct result table
    DROP TEMPORARY TABLE IF EXISTS _payoutList;
    
    CREATE TEMPORARY TABLE _payoutList (
        payoutID       INT,
        payDate        DATE,
        regionCode     VARCHAR(10),
        fixedBonus     DECIMAL(20,2),
        variableBonus  DECIMAL(20,2),
        transportation DECIMAL(20,2),
        addon          DECIMAL(20,2),
        deduction      DECIMAL(20,2),
        total          DECIMAL(20,2),
        totalBudget    DECIMAL(20,2),
        difference     DECIMAL(20,2) GENERATED ALWAYS AS (totalBudget - total),
        isPaid         BOOLEAN,
        status         BOOLEAN,
        totalCount     INT
    );

    -- Query ungrouped payout data
    IF _isGrouped IS NULL OR _isGrouped IS FALSE THEN
        IF _isPaid IS NULL OR _isPaid IS FALSE THEN
            -- Count records
            SELECT COUNT(DISTINCT C.payDate, C.regionCode) INTO _ungroupedCount
            FROM Campaign C
            INNER JOIN CampaignLeaderboard L ON C.id = L.campaignID
            INNER JOIN CampaignRank R ON L.id = R.leaderboardID
            INNER JOIN CampaignBonus B ON R.id = B.rankID AND (B.payDate IS NULL OR B.payDate = C.payDate)
            LEFT JOIN PayoutDetail D ON B.id = D.bonusID
            WHERE C.approvalStatus IS TRUE
              AND D.payoutID IS NULL
              AND (_payDate IS NULL OR C.payDate = _payDate)
              AND (_regionList IS NULL OR FIND_IN_SET(C.regionCode, _regionList) > 0);
            
            -- Place current page data
            INSERT INTO _payoutList (payoutID, payDate, regionCode, isPaid, status)
            SELECT 0, C.payDate, C.regionCode, FALSE, 0
            FROM Campaign C
            INNER JOIN CampaignLeaderboard L ON C.id = L.campaignID
            INNER JOIN CampaignRank R ON L.id = R.leaderboardID
            INNER JOIN CampaignBonus B ON R.id = B.rankID AND (B.payDate IS NULL OR B.payDate = C.payDate)
            LEFT JOIN PayoutDetail D ON B.id = D.bonusID
            WHERE C.approvalStatus IS TRUE
              AND D.payoutID IS NULL
              AND (_payDate IS NULL OR C.payDate = _payDate)
              AND (_regionList IS NULL OR FIND_IN_SET(C.regionCode, _regionList) > 0)
            GROUP BY C.payDate, C.regionCode
            ORDER BY payDate, regionCode
            LIMIT _limitCount OFFSET _offsetCount;
            
            -- Update bonus summary data
            UPDATE _payoutList P
            INNER JOIN (
                SELECT payDate, regionCode, SUM(IF(bonusType = 0, amount, 0)) AS fixedBonus, 
                    SUM(IF(bonusType = 1, amount, 0)) AS variableBonus,
                    SUM(IF(bonusType = 2, amount, 0)) AS transportation,
                    SUM(IF(bonusType = 3, amount, 0)) AS addon,
                    SUM(IF(bonusType = 4, amount, 0)) AS deduction,
                    SUM(amount) AS total
                FROM v_payoutDetail2
                GROUP BY regionCode, payDate
            ) S ON P.payDate = S.payDate AND P.regionCode = S.regionCode
            SET P.fixedBonus = S.fixedBonus, P.variableBonus = S.variableBonus, P.transportation = S.transportation, P.addon = S.addon, P.deduction = S.deduction, P.total = S.total;
    
            -- Update record count
            SELECT COUNT(*) INTO _pagedCount
            FROM _payoutList;
        END IF;
    END IF;
    
    -- Query grouped payout data
    IF _isGrouped IS NULL OR _isGrouped IS TRUE THEN
        -- Count records
        SELECT COUNT(*) INTO _groupedCount
        FROM Payout P
        WHERE (_payDate IS NULL OR P.payDate = _payDate)
          AND (_regionList IS NULL OR FIND_IN_SET(P.regionCode, _regionList) > 0);
        
        IF _offsetCount > _ungroupedCount OR _limitCount > _pagedCount THEN
            -- Adjust offset and limit
            SET _offsetCount = IF(_offsetCount > _ungroupedCount, _offsetCount - _ungroupedCount, _offsetCount);
            SET _limitCount = IF(_limitCount > _pagedCount, _limitCount - _pagedCount, _limitCount);
            
            -- Place current page data
            INSERT INTO _payoutList (payoutID, payDate, regionCode, isPaid, status)
            SELECT id, payDate, regionCode, payStatus, IF(payStatus, 2, 1) AS status
            FROM Payout
            WHERE (_payDate IS NULL OR payDate = _payDate)
              AND (_regionList IS NULL OR FIND_IN_SET(regionCode, _regionList) > 0)
            ORDER BY payDate, regionCode
            LIMIT _limitCount OFFSET _offsetCount;
            
            -- Update bonus summary data
            UPDATE _payoutList L
            INNER JOIN (
                SELECT D.payoutID, SUM(IF(B.bonusType = 0, B.amount, 0)) AS fixedBonus, 
                    SUM(IF(B.bonusType = 1, B.amount, 0)) AS variableBonus,
                    SUM(IF(B.bonusType = 2, B.amount, 0)) AS transportation,
                    SUM(IF(B.bonusType = 3, B.amount, 0)) AS addon,
                    SUM(IF(B.bonusType = 4, B.amount, 0)) AS deduction,
                    SUM(B.amount) AS total
                FROM CampaignBonus B
                INNER JOIN PayoutDetail D ON B.id = D.bonusID
                GROUP BY D.payoutID
            ) S ON L.payoutID = S.payoutID
            SET L.fixedBonus = S.fixedBonus, L.variableBonus = S.variableBonus, L.transportation = S.transportation, L.addon = S.addon, L.deduction = S.deduction, L.total = S.total
            WHERE L.payoutID > 0;            
        END IF;
    END IF;
    
    -- Update campaign totalBudget data
    UPDATE _payoutList P
    INNER JOIN (
        SELECT regionCode, payDate, SUM(budget) AS totalBudget
        FROM Campaign
        WHERE approvalStatus IS TRUE
        GROUP BY regionCode, payDate
    ) C ON P.payDate = C.payDate AND P.regionCode = C.regionCode
    SET P.totalBudget = C.totalBudget;
    
    -- Place totalCount
    UPDATE _payoutList
    SET totalCount = (_ungroupedCount + _groupedCount)
    LIMIT 1;
    
    -- Output result data
    SELECT *
    FROM _payoutList
    ORDER BY payoutID = 0 DESC, payoutID DESC, payDate, regionCode;

    -- Release temp table
    DROP TEMPORARY TABLE _payoutList;
END //

DELIMITER ;
