package service

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/mailServer"
	"17live_wso_be/internal/model"
	"17live_wso_be/util"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	existedRegions map[string]string
	newRegions     []model.Region

	noRegionCampaign []int
)

func (c *Client) SyncData(ctx context.Context, userEmail string) error {
	log.Infof("service sync data")

	startTime := time.Now().UTC()

	skipList, err := c.RepositoryClient.GetSkipSyncList(ctx)
	if err != nil {
		return err
	}

	campaignApiData, err := c.MediaClient.FetchCampaignData(ctx, skipList)
	if err != nil {
		log.Errorf("fail to fetch 17 campaign data: %s", err.Error())
		return customError.New(customError.FetchCampaignFail)
	}

	regions, err := c.GetAllRegionCodeIncludingUnshowable(ctx)
	if err != nil {
		return err
	}
	existedRegions = util.StringSliceToMap(regions)

	campaignSet, err := c.groupCampaignSet(ctx, campaignApiData)
	if err != nil {
		return err
	}

	campaignSet.Regions = newRegions

	log.Infof("ready to upsert %d campaign, %d leaderboard, %d streamerAgency, %d streamer, %d rank", len(campaignSet.Campaigns), len(campaignSet.Leaderboards), len(campaignSet.StreamerAgencies), len(campaignSet.Streamers), len(campaignSet.Ranks))

	if err := c.RepositoryClient.UpsertCampaignRelatedData(ctx, *campaignSet.SplitData()); err != nil {
		return err
	}

	if len(noRegionCampaign) == 0 {
		log.Infof("synced campaign all got region, skip sending notification email")
		return nil
	}

	mailList, err := c.GetAdminMailList(ctx)
	if err != nil {
		return err
	}

	if err := mailServer.New().SendNoRegionNotification(ctx, mailList, noRegionCampaign); err != nil {
		log.Errorf("campaign related data synced, but fail to send no region notification email: %s", err)
		return customError.New(customError.MailDeliveryFailed)
	}

	if err := mailServer.New().SendSyncDataFinishNotification(ctx, userEmail); err != nil {
		log.Errorf("campaign related data synced, but fail to send sync data finish notification email: %s", err)
		return customError.New(customError.MailDeliveryFailed)
	}

	log.Infof("finish sync data, time duration: %v", time.Since(startTime))

	return nil
}

func (c *Client) groupCampaignSet(ctx context.Context, campaignApiData []model.MediaApiCampaignSet) (model.CampaignSet, error) {
	log.Infof("start group campaign set, raw campaign count: %d", len(campaignApiData))

	time := time.Now().UTC()

	var data model.CampaignSet

	streamerIdMap := make(map[uuid.UUID]string)
	accountIdMap := make(map[int]string)

	for _, d := range campaignApiData {
		var campaign model.CampaignInsert

		var region *string
		var regionList *string
		if len(d.Campaign.Regions) != 0 && d.Campaign.Regions[0] != "GLOBAL" {
			for _, region := range d.Campaign.Regions {
				c.checkRegion(ctx, region)
			}

			r := strings.Join(d.Campaign.Regions, ",")
			regionList = &r

			if len(d.Campaign.Regions) == 1 {
				regionList = nil
				region = &d.Campaign.Regions[0]
			}
		}

		campaign = model.CampaignInsert{
			Id:         d.Campaign.Id,
			Title:      getCampaignTitle(ctx, d),
			RegionCode: region,
			RegionList: regionList,
			StartTime:  d.Campaign.StartTime,
			EndTime:    d.Campaign.EndTime,
			SyncTime:   time,
		}

		token, err := c.MediaClient.FetchAccessToken(ctx)
		if err != nil {
			log.Errorf("fail to fetch access token from 17 media: %s", err.Error())
			return data, customError.New(customError.FetchAccessTokenFail)
		}

		leaderboardCount := 0

		for _, l := range d.Leaderboards {
			leaderboard := model.CampaignLeaderboard{
				Id:         l.Id,
				CampaignID: d.Campaign.Id,
				Title:      l.GroupName,
				SyncTime:   time,
			}

			leaderboardDetailApiData, err := c.MediaClient.FetchLeaderboardData(ctx, l.Id.String())
			if err != nil {
				log.Errorf("fail to fetch 17 leaderboard rank of campaign %d, leaderborad %s: %s", d.Campaign.Id, l.Id.String(), err.Error())
				// use continue to avoid interrupting by dirty data
				continue
				// return data, customError.New(customError.FetchLeaderboardFail)
			}

			rankCount := 0

			for _, r := range leaderboardDetailApiData.Ranks {
				if _, ok := streamerIdMap[r.Streamer]; !ok {
					streamer := model.Streamer{
						Id:       r.Streamer,
						SyncTime: time,
					}

					streamer, streamerAgency, err := c.groupStreamerData(ctx, streamer, token, accountIdMap)
					if err != nil {
						return data, err
					}

					if streamerAgency.Id != 0 {
						data.StreamerAgencies = append(data.StreamerAgencies, streamerAgency)
					}

					data.Streamers = append(data.Streamers, *streamer.UUIDToBinary())
					streamerIdMap[r.Streamer] = ""
				}

				rank := model.CampaignRank{
					LeaderboardID:  l.Id,
					Rank:           r.Rank,
					Score:          r.Score,
					StreamerUserID: r.Streamer,
					SyncTime:       time,
				}

				data.Ranks = append(data.Ranks, *rank.UUIDToBinary())
				rankCount += 1
			}

			if rankCount == 0 {
				// log.Infof("no rank, skip leaderboard %s of campaign %d", l.Id.String(), campaign.Id)
				continue
			}

			data.Leaderboards = append(data.Leaderboards, *leaderboard.UUIDToBinary())
			leaderboardCount += 1
		}

		if leaderboardCount == 0 {
			log.Infof("no leaderboard, skip campaign %d", campaign.Id)
			continue
		}

		if campaign.RegionCode == nil {
			noRegionCampaign = append(noRegionCampaign, campaign.Id)
		}

		data.Campaigns = append(data.Campaigns, campaign)
	}

	return data, nil
}

func (c *Client) groupStreamerData(ctx context.Context, streamer model.Streamer, token string, accountIdMap map[int]string) (model.Streamer, model.StreamerAgency, error) {
	var streamerAgency model.StreamerAgency

	streamerApiData, err := c.MediaClient.FetchStreamerData(ctx, streamer.Id.String(), token)
	if err != nil {
		// log.Errorf("fail to fetch 17 streamer info of %s: %s", streamer.Id.String(), err.Error())
		return streamer, streamerAgency, nil // return nil to avoid interrupting by dirty data
	} else {
		streamer.Name = &streamerApiData.Name
		streamer.OpenID = &streamerApiData.OpenID
	}

	streamerContractApiData, err := c.MediaClient.FetchStreamerContractData(ctx, streamer.Id.String(), token)
	if err != nil {
		// log.Errorf("fail to fetch 17 streamer contract of %s: %s", streamer.Id.String(), err.Error())
		return streamer, streamerAgency, nil // return nil to avoid interrupting by dirty data
	}

	accountData := streamerContractApiData.PayoutAccounts.Third

	accountID, err := c.getBankAccount(ctx, streamer.Id, accountData, streamer.SyncTime, accountIdMap)
	if err != nil {
		return streamer, streamerAgency, err
	}

	if streamerContractApiData.Region != "" && streamerContractApiData.Region != "GLOBAL" {
		c.checkRegion(ctx, streamerContractApiData.Region)
		streamer.RegionCode = &streamerContractApiData.Region
	}

	if streamerContractApiData.AgencyID != nil {
		streamerAgency.Id = *streamerContractApiData.AgencyID
		streamerAgency.Name = streamerContractApiData.AgencyName
		streamerAgency.SyncTime = streamer.SyncTime

		if accountData.OwnerType != nil && *accountData.OwnerType == model.OwnerTypeAgency {
			streamerAgency.AccountID = accountID
			streamer.AgencyID = streamerContractApiData.AgencyID
		}
	}

	if accountData.OwnerType != nil && *accountData.OwnerType == model.OwnerTypeStreamer {
		streamer.AccountID = accountID
	}

	return streamer, streamerAgency, nil
}

func (c *Client) getBankAccount(ctx context.Context, streamer uuid.UUID, accountData model.MediaApiStreamerAccount, syncTime time.Time, accountIdMap map[int]string) (*int, error) {
	var accountID *int

	if accountData.AccountType == nil || accountData.BankAccount == "" || accountData.PayeeName == nil || accountData.PayeeType == nil {
		return accountID, nil
	}

	if accountData.PayeeID != nil && strings.Count(*accountData.PayeeID, "") > 11 {
		log.Errorf("fail to save streamer %s account data, payeeID too long: %s", streamer.String(), *accountData.PayeeID)
		return accountID, nil
	}

	var feeOwnerType *int

	if accountData.FeeOwnerType != 0 {
		feeOwnerType = &accountData.FeeOwnerType
	}

	account := model.BankAccount{
		SwiftCode:       accountData.SwiftCode,
		BankCode:        accountData.BankCode,
		BankName:        accountData.BankName,
		BranchCode:      accountData.BranchCode,
		BranchName:      accountData.BranchName,
		BankAccount:     accountData.BankAccount,
		BankAccountName: *accountData.PayeeName,
		BankAccountType: *accountData.AccountType,
		PayeeType:       *accountData.PayeeType,
		PayeeID:         accountData.PayeeID,
		OwnerType:       accountData.OwnerType,
		FeeOwnerType:    feeOwnerType,
		SyncTime:        syncTime,
	}

	res, err := c.GetBankAccount(ctx, account.SwiftCode, account.BankCode, account.BranchCode, account.BankAccount)
	if err != nil {
		return accountID, err
	} else if len(res) == 0 {
		if created, err := c.CreateBankAccount(ctx, account); err != nil {
			return accountID, err
		} else {
			accountID = &created.Id
			accountIdMap[*accountID] = ""
			return accountID, nil
		}
	}

	account.Id = res[0].Id

	if _, ok := accountIdMap[account.Id]; !ok {
		if err := c.UpdateBankAccount(ctx, account); err != nil {
			return accountID, err
		}
		accountIdMap[account.Id] = ""
	}

	accountID = &account.Id

	return accountID, nil
}

func (c *Client) checkRegion(ctx context.Context, region string) {
	if _, ok := existedRegions[region]; !ok {
		data := model.Region{
			Code: region,
			Name: region,
		}

		newRegions = append(newRegions, data)
	}
}

func getCampaignTitle(ctx context.Context, data model.MediaApiCampaignSet) string {
	var title string

	switch data.Campaign.DefaultLang {
	case model.CampaignLangAR:
		title = data.Campaign.Names.AR
	case model.CampaignLangEN:
		title = data.Campaign.Names.EN
	case model.CampaignLangUS:
		title = data.Campaign.Names.US
	case model.CampaignLangJP:
		title = data.Campaign.Names.JP
	case model.CampaignLangCN:
		title = data.Campaign.Names.CN
	case model.CampaignLangHK:
		title = data.Campaign.Names.HK
	case model.CampaignLangTW:
		title = data.Campaign.Names.TW
	}

	if title == "" {
		title = data.Campaign.OpenName
	}

	return title
}
