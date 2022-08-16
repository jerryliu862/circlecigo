package repository

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"context"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func (c *Client) ListCampaign(ctx context.Context, filter model.CampaignFilter) ([]model.Campaign, int, error) {
	var data []model.Campaign
	var total int

	filterSQL := parseCampaignFilterSQL(ctx, filter)

	if err := c.Database.Raw("SELECT c.*, c.budget - c.totalBonus AS bonusDiff, (SELECT COUNT(*) FROM v_campaignList c WHERE c.regionCode IN (?)"+filterSQL+") AS totalCount FROM v_campaignList c WHERE c.regionCode IN (?)"+filterSQL+" ORDER BY id LIMIT ? OFFSET ?", filter.Regions, filter.Regions, filter.PageSize, filter.PageSize*(filter.PageNum-1)).Find(&data).Error; err != nil {
		log.Errorf("fail to list campaign with filter %v: %s", filter, err.Error())
		return data, total, customError.New(customError.DatabaseError)
	}

	if len(data) != 0 {
		total = data[0].TotalCount
	}

	return data, total, nil
}

func (c *Client) ListNoRegionCampaign(ctx context.Context, page model.PageFilter) ([]model.CampaignBasic, int, error) {
	var data []model.CampaignBasic
	var total int

	if err := c.Database.Raw("SELECT id, title, startTime, endTime, regionList, (SELECT COUNT(*) FROM Campaign WHERE regionCode IS NULL) AS totalCount FROM Campaign WHERE regionCode IS NULL ORDER BY id LIMIT ? OFFSET ?", page.PageSize, page.PageSize*(page.PageNum-1)).Find(&data).Error; err != nil {
		log.Errorf("fail to list no region campaign with filter %v: %s", page, err.Error())
		return data, total, customError.New(customError.DatabaseError)
	}

	if len(data) != 0 {
		total = data[0].TotalCount
	}

	return data, total, nil
}

func (c *Client) ListUnpaidCampaignBonus(ctx context.Context, regions []string) ([]model.CampaignBonusUnpaid, error) {
	var data []model.CampaignBonusUnpaid

	if err := c.Database.Table(ViewUnpaidBonusList).Select("rankID", "userID", "openID", "campaignID", "campaignTitle", "campaignCurrency").Where("regionCode IN (?)", regions).Find(&data).Error; err != nil {
		log.Errorf("fail to list unpaid campaign bonus of streamer regions: %v, %s", regions, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetCampaignById(ctx context.Context, id int) ([]model.Campaign, error) {
	var data []model.Campaign

	if err := c.Database.Table(Campaign).Where("id = ?", id).Find(&data).Error; err != nil {
		log.Errorf("fail to get campaign: %d. %s", id, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetCampaignByRankId(ctx context.Context, id int) (model.Campaign, error) {
	var data model.Campaign

	if err := c.Database.Raw("SELECT c.* FROM Campaign c LEFT JOIN CampaignLeaderboard cl on cl.campaignID = c.id LEFT JOIN CampaignRank cr on cr.leaderboardID = cl.id WHERE cr.id = ?", id).First(&data).Error; err != nil {
		log.Errorf("fail to get campaign by rank id: %d. %s", id, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetCampaignByBonusId(ctx context.Context, id int) (model.Campaign, error) {
	var data model.Campaign

	if err := c.Database.Raw("SELECT c.* FROM Campaign c LEFT JOIN CampaignLeaderboard cl on cl.campaignID = c.id LEFT JOIN CampaignRank cr on cr.leaderboardID = cl.id LEFT JOIN CampaignBonus cb on cb.rankID = cr.id WHERE cb.id = ?", id).First(&data).Error; err != nil {
		log.Errorf("fail to get campaign by bonus id: %d. %s", id, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetCampaignLeaderboard(ctx context.Context, leaderboardID []byte) (model.CampaignLeaderboard, error) {
	var data model.CampaignLeaderboard

	if err := c.Database.Table(CampaignLeaderboard).Where("id = ?", leaderboardID).First(&data).Error; err != nil {
		log.Errorf("fail to get campaign leaderboard: %s", err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetCampaignRank(ctx context.Context, rankID int) (model.CampaignRank, error) {
	var data model.CampaignRank

	if err := c.Database.Table(CampaignRank).Where("id = ?", rankID).First(&data).Error; err != nil {
		log.Errorf("fail to get campaign rank: %s", err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	data = *data.BinaryToUUID()

	return data, nil
}

func (c *Client) GetCampaignBonusById(ctx context.Context, id int) ([]model.CampaignBonus, error) {
	var data []model.CampaignBonus

	if err := c.Database.Table(CampaignBonus).Where("id = ?", id).Find(&data).Error; err != nil {
		log.Errorf("fail to get campaign bonus by bonus id: %d. %s", id, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetCampaignBonusByRankIdAndBonusType(ctx context.Context, rankID int, bonusType int) ([]model.CampaignBonus, error) {
	var data []model.CampaignBonus

	if err := c.Database.Table(CampaignBonus).Where("rankID = ? AND bonusType = ?", rankID, bonusType).Find(&data).Error; err != nil {
		log.Errorf("fail to get campaign bonus: rank %d, bonusType %d. %s", rankID, bonusType, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetCampaignApprovalStatus(ctx context.Context, id int) (bool, error) {
	var status bool

	if err := c.Database.Table(Campaign).Model(&model.Campaign{}).Select("approvalStatus").Where("id = ?", id).First(&status).Error; err != nil {
		log.Errorf("fail to get campaign approval status: %d. %s", id, err.Error())
		return status, customError.New(customError.DatabaseError)
	}

	return status, nil
}

func (c *Client) GetCampaignDetail(ctx context.Context, id int) (model.CampaignDetail, error) {
	var data model.CampaignDetail

	err := c.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(ViewCampaignList).Where("id = ?", id).Select("*, budget - totalBonus AS bonusDiff").Find(&data.Campaign).Error; err != nil {
			return err
		}

		if err := tx.Raw("CALL p_GetCampaignDetailByStage(?)", id).Scan(&data.LeaderboardList).Error; err != nil {
			return err
		}

		for i := range data.LeaderboardList {
			data.LeaderboardList[i] = *data.LeaderboardList[i].BinaryToUUID()
			if err := tx.Raw("CALL p_GetCampaignDetailByRank(UUID_TO_BIN(?))", data.LeaderboardList[i].LeaderboardID).Scan(&data.LeaderboardList[i].RankList).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Errorf("fail to get campaign datail: %d. %s", id, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) ApproveCampaignBonus(ctx context.Context, id int, approver int) error {
	if err := c.Database.Exec("UPDATE Campaign SET approvalStatus = ?, approvalTime = ?, approverUID = ? WHERE id = ?", true, time.Now().UTC(), approver, id).Error; err != nil {
		log.Errorf("fail to approve campaign bonus of %d: %s", id, err.Error())
		return customError.New(customError.DatabaseError)
	}

	return nil
}

func (c *Client) UpdateCampaignRegion(ctx context.Context, id int, region string, modifier int) error {
	if err := c.Database.Exec("UPDATE Campaign SET regionCode = ?, modifyTime = ?, modifierUID = ? WHERE id = ?", region, time.Now().UTC(), modifier, id).Error; err != nil {
		log.Errorf("fail to update campaign region of %d: %s", id, err.Error())
		return customError.New(customError.DatabaseError)
	}

	return nil
}

func (c *Client) SetCampaignBonusPaydate(ctx context.Context, date string, bonusIdList []int, modifier int) error {
	if err := c.Database.Exec("UPDATE CampaignBonus SET payDate = ?, modifyTime = ?, modifierUID = ? WHERE id IN (?)", date, time.Now().UTC(), modifier, bonusIdList).Error; err != nil {
		log.Errorf("fail to set campaign bonus date as: %s, bonusIdList %v, %s", date, bonusIdList, err.Error())
		return customError.New(customError.DatabaseError)
	}

	return nil
}

func (c *Client) UpsertCampaignAndBonus(ctx context.Context, data model.CampaignDetail) error {
	log.Infof("upsert campaign and bonus: %v", data)

	err := c.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE Campaign SET budget = ?, remark = ?, modifyTime = ?, modifierUID = ? WHERE id = ?", data.Budget, data.Remark, time.Now().UTC(), data.ModifierUID, data.Id).Error; err != nil {
			return err
		}

		if len(data.BonusList) != 0 {
			for i, bonus := range data.BonusList {
				if err := tx.Table(CampaignBonus).Where("rankID = ? AND bonusType = ?", bonus.RankID, bonus.BonusType).Select("id").Find(&data.BonusList[i].Id).Error; err != nil {
					return err
				}

				if data.BonusList[i].Id != 0 {
					if err := tx.Exec("UPDATE CampaignBonus SET amount = ?, modifyTime = ?, modifierUID = ? WHERE id = ?", bonus.Amount, bonus.ModifyTime, bonus.ModifierUID, data.BonusList[i].Id).Error; err != nil {
						return err
					}
				} else {
					if err := tx.Table(CampaignBonus).Create(&bonus).Error; err != nil {
						return err
					}
				}
			}

			var campaign model.Campaign
			if err := tx.Raw("SELECT SUM(cb.amount) AS totalBonus FROM CampaignBonus cb LEFT JOIN CampaignRank cr on cb.rankID = cr.id LEFT JOIN CampaignLeaderboard cl on cl.id = cr.leaderboardID LEFT JOIN Campaign c on c.id = cl.campaignID WHERE c.id = ? AND cb.bonusType IN (?)", data.Id, []int{model.BonusTypeFixed, model.BonusTypeVariable}).Find(&campaign).Error; err != nil {
				return err
			}

			if err := tx.Exec("UPDATE Campaign SET totalBonus = ? WHERE id = ?", campaign.TotalBonus, data.Id).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Errorf("fail to upsert campaign and bonus: %s", err.Error())
		return customError.New(customError.DatabaseError)
	}

	log.Infof("campaign related data upserted")

	return nil
}

func (c *Client) CreateCampaignBonus(ctx context.Context, data model.CampaignBonus) (model.CampaignBonus, error) {
	if err := c.Database.Table(CampaignBonus).Create(&data).Error; err != nil {
		log.Errorf("fail to create campaign bonus: %v. %s", data, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	log.Infof("campaign bonus created: %v", data)

	return data, nil
}

func (c *Client) UpdateCampaignBonus(ctx context.Context, data model.CampaignBonus) error {
	if err := c.Database.Table(CampaignBonus).Save(&data).Error; err != nil {
		log.Errorf("fail to update campaign bonus: %v. %s", data, err.Error())
		return customError.New(customError.DatabaseError)
	}

	log.Infof("campaign bonus updated: %v", data)

	return nil
}

func (c *Client) DeleteCampaignBonus(ctx context.Context, id int) error {
	if err := c.Database.Table(CampaignBonus).Where("id = ?", id).Delete(&model.CampaignBonus{}).Error; err != nil {
		log.Errorf("fail to delete campaign bonus: %d. %s", id, err.Error())
		return customError.New(customError.DatabaseError)
	}

	log.Infof("campaign bonus deleted: %d", id)

	return nil
}

func parseCampaignFilterSQL(ctx context.Context, filter model.CampaignFilter) string {
	var res string

	if filter.Date == nil && filter.Approval == nil && filter.IsZero == nil && filter.Keyword == "" {
		return res
	}

	if filter.Date != nil {
		res += fmt.Sprintf(" AND YEAR(FROM_UNIXTIME(c.endTime)) = %d AND MONTH(FROM_UNIXTIME(c.endTime)) = %d", filter.Date.Year(), int(filter.Date.Month()))
	}

	if filter.Approval != nil {
		res += fmt.Sprintf(" AND c.approvalStatus = %t", *filter.Approval)
	}

	if filter.IsZero != nil {
		if *filter.IsZero {
			res += " AND c.totalBonus = 0"
		} else {
			res += " AND c.totalBonus <> 0"
		}
	}

	if filter.Keyword != "" {
		if id, err := strconv.Atoi(filter.Keyword); err != nil {
			res += fmt.Sprintf(" AND c.title LIKE '%%%s%%'", filter.Keyword)
		} else {
			res += fmt.Sprintf(" AND (c.title LIKE '%%%s%%' OR c.id = %d)", filter.Keyword, id)
		}
	}

	log.Infof("parse campaign filter SQL with filter %v as: %s", filter, res)

	return res
}
