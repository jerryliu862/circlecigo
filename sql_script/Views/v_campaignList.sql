CREATE OR REPLACE VIEW v_campaignList
AS
    SELECT DISTINCT C.id, C.title, C.regionCode, C.startTime, C2.endTime, C2.budget, C2.totalBonus, C2.budget - C2.totalBonus AS bonusDiff, C.remark, C.syncTime, C.approvalStatus, C.approvalTime, C2.approverName
    FROM Campaign C
    INNER JOIN (
        SELECT id, endTime -1 AS endTime, IFNULL(budget, 0) AS budget, fn_GetCampaignTotalBonus(id) AS totalBonus, fn_GetUserName(approverUID) AS approverName
        FROM Campaign
    ) C2 ON C.id = C2.id
