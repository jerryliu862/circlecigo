package service

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"17live_wso_be/util"
	"context"
	"strings"
	"time"
)

func (c *Client) ListCampaign(ctx context.Context, filter model.CampaignFilter) ([]model.Campaign, int, error) {
	log.Infof("service list campaign with filter: %v", filter)
	return c.RepositoryClient.ListCampaign(ctx, filter)
}

func (c *Client) ListNoRegionCampaign(ctx context.Context, page model.PageFilter) ([]model.CampaignBasic, int, error) {
	log.Infof("service list no region campaign with filter: %v", page)
	return c.RepositoryClient.ListNoRegionCampaign(ctx, page)
}

func (c *Client) GetCampaignById(ctx context.Context, id int) ([]model.Campaign, error) {
	log.Infof("service get campaign by id: %d", id)
	return c.RepositoryClient.GetCampaignById(ctx, id)
}

func (c *Client) GetCampaignDetail(ctx context.Context, id int) (model.CampaignDetail, error) {
	log.Infof("service get campaign detail: %d", id)
	return c.RepositoryClient.GetCampaignDetail(ctx, id)
}

func (c *Client) SetCampaignBonus(ctx context.Context, uid int, data model.CampaignDetail) error {
	log.Infof("service set campaign bonus: %v", data)

	if approved, err := c.RepositoryClient.GetCampaignApprovalStatus(ctx, data.Id); err != nil {
		return err
	} else if approved {
		log.Warnf("campaign bonus has approved: campaignID %d", data.Id)
		return customError.New(customError.CampaignBonusApproved)
	}

	data.ModifierUID = &uid

	for _, leaderboard := range data.LeaderboardList {
		if l, err := c.RepositoryClient.GetCampaignLeaderboard(ctx, leaderboard.UUIDToBinary().BinaryLeaderboardID); err != nil || l.CampaignID != data.Id {
			log.Warnf("fail to get campaign leaderboard: campaignID %d, leaderboardID %s", data.Id, leaderboard.LeaderboardID.String())
			return customError.New(customError.LeaderboardNotInCampaign)
		}

		for _, rank := range leaderboard.RankList {
			if r, err := c.RepositoryClient.GetCampaignRank(ctx, rank.RankID); err != nil || r.LeaderboardID != leaderboard.LeaderboardID {
				log.Warnf("fail to get campaign rank: campaignID %d, leaderboardID %s, rankID %d", data.Id, leaderboard.LeaderboardID.String(), rank.RankID)
				return customError.New(customError.RankNotInLeaderboard)
			}
			data.BonusList = append(data.BonusList, parseCampaignBonus(ctx, uid, rank)...)
		}
	}

	return c.RepositoryClient.UpsertCampaignAndBonus(ctx, data)
}

func (c *Client) ApproveCampaignBonus(ctx context.Context, id int, approver int) error {
	log.Infof("service approve campaign bonus: %d, approver %d", id, approver)
	return c.RepositoryClient.ApproveCampaignBonus(ctx, id, approver)
}

func (c *Client) ListUnpaidCampaignBonus(ctx context.Context, uid int) ([]model.CampaignBonusUnpaid, error) {
	log.Infof("service list unpaid campaign bonus for payout admin: uid %d", uid)

	regions, err := c.GetUserAuthRegion(ctx, uid, model.UserAuthTypePayout, model.UserAuthLevelEdit)
	if err != nil {
		return nil, err
	}

	return c.RepositoryClient.ListUnpaidCampaignBonus(ctx, regions)
}

func (c *Client) SetCampaignRegion(ctx context.Context, data []model.Campaign, uid int) error {
	log.Info("service set campaign region")

	for _, campaign := range data {
		d, err := c.GetCampaignById(ctx, campaign.Id)
		if err != nil {
			return err
		} else if len(d) == 0 {
			log.Warnf("check campaign before set region, campaign does not exist, request campaign id: %d", campaign.Id)
			return customError.New(customError.CampaignNotExist)
		} else if d[0].RegionCode != nil {
			log.Warnf("check campaign before set region, campaign %d has set region as: %s", campaign.Id, *d[0].RegionCode)
			return customError.New(customError.CampaignRegionSet)
		}

		if campaign.RegionCode == nil {
			log.Warnf("empty request region for campaign %d", campaign.Id)
			return customError.New(customError.InvalidRequestData)
		}

		region := strings.ToUpper(*campaign.RegionCode)
		if _, err := c.GetRegion(ctx, region); err != nil {
			log.Warnf("region set for campaign %d does not exist: %s", campaign.Id, region)
			return customError.New(customError.InvalidRequestData)
		}

		if d[0].RegionList != nil {
			regionList := strings.Split(*d[0].RegionList, ",")
			if !util.ContainString(regionList, region) {
				log.Warnf("region set for campaign %d does not in its region list: %s, region list %v", campaign.Id, region, regionList)
				return customError.New(customError.CampaignRegionNotInList)
			}
		}

		c.RepositoryClient.UpdateCampaignRegion(ctx, campaign.Id, region, uid)
	}

	return nil
}

func (c *Client) CheckCampaignBonusExistence(ctx context.Context, ids []int) error {
	for _, id := range ids {
		if b, err := c.RepositoryClient.GetCampaignBonusById(ctx, id); err != nil {
			return err
		} else if len(b) == 0 {
			log.Warnf("campaign bonus does not exist: %d", id)
			return customError.New(customError.CampaignBonusNotExist)
		}
	}

	return nil
}

func (c *Client) CheckCampaignBonusPermission(ctx context.Context, uid int, authType int, authLevel int, ids []int) bool {
	for _, id := range ids {
		campaign, err := c.RepositoryClient.GetCampaignByBonusId(ctx, id)
		if err != nil {
			return false
		} else if campaign.RegionCode == nil {
			log.Warnf("region code of campaign %d is null", campaign.Id)
			return false
		}

		if !c.PermissionCheck(ctx, uid, *campaign.RegionCode, authType, authLevel) {
			return false
		}
	}

	return true
}

func (c *Client) CheckCampaignRankPermission(ctx context.Context, uid int, authType int, authLevel int, id int) bool {
	campaign, err := c.RepositoryClient.GetCampaignByRankId(ctx, id)
	if err != nil {
		return false
	} else if campaign.RegionCode == nil {
		log.Warnf("region code of campaign %d is null", campaign.Id)
		return false
	}

	if !c.PermissionCheck(ctx, uid, *campaign.RegionCode, authType, authLevel) {
		return false
	}

	return true
}

func parseCampaignBonus(ctx context.Context, uid int, rank model.RankDetail) []model.CampaignBonus {
	var bonusList []model.CampaignBonus

	time := time.Now().UTC()

	fixedBonus := model.CampaignBonus{
		RankID:      rank.RankID,
		BonusType:   model.BonusTypeFixed,
		Amount:      rank.FixedBonus,
		CreateTime:  time,
		CreatorUID:  uid,
		ModifyTime:  time,
		ModifierUID: uid,
	}

	variableBonus := model.CampaignBonus{
		RankID:      rank.RankID,
		BonusType:   model.BonusTypeVariable,
		Amount:      rank.VariableBonus,
		CreateTime:  time,
		CreatorUID:  uid,
		ModifyTime:  time,
		ModifierUID: uid,
	}

	bonusList = append(bonusList, fixedBonus)
	bonusList = append(bonusList, variableBonus)

	return bonusList
}
