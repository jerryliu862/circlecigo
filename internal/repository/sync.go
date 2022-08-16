package repository

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *Client) UpsertCampaignRelatedData(ctx context.Context, data model.CampaignSet) error {
	err := c.Database.Transaction(func(tx *gorm.DB) error {
		if len(data.Regions) != 0 {
			if err := tx.Table(Region).Clauses(clause.OnConflict{DoNothing: true}).Create(&data.Regions).Error; err != nil {
				return err
			}
		}

		for _, d := range data.CampaignsDataSplit {
			if len(d) != 0 {
				if err := tx.Table(Campaign).Clauses(clause.OnConflict{
					DoUpdates: clause.AssignmentColumns([]string{"title", "regionCode", "regionList", "startTime", "endTime", "syncTime"}),
				}).Create(&d).Error; err != nil {
					return err
				}
			}
		}

		for _, d := range data.LeaderboardsDataSplit {
			if len(d) != 0 {
				if err := tx.Table(CampaignLeaderboard).Clauses(clause.OnConflict{UpdateAll: true}).Create(&d).Error; err != nil {
					return err
				}
			}
		}

		for _, d := range data.StreamerAgenciesDataSplit {
			if len(d) != 0 {
				if err := tx.Table(StreamerAgency).Clauses(clause.OnConflict{UpdateAll: true}).Create(&d).Error; err != nil {
					return err
				}
			}
		}

		for _, d := range data.StreamersDataSplit {
			if len(d) != 0 {
				if err := tx.Table(Streamer).Clauses(clause.OnConflict{UpdateAll: true}).Create(&d).Error; err != nil {
					return err
				}
			}
		}

		for _, d := range data.RanksDataSplit {
			if len(d) != 0 {
				if err := tx.Table(CampaignRank).Clauses(clause.OnConflict{UpdateAll: true}).Create(&d).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Errorf("fail to upsert campaign related data: %s", err.Error())
		return customError.New(customError.DatabaseError)
	}

	// log.Info("campaign related data upserted")

	return nil
}

func (c *Client) GetSkipSyncList(ctx context.Context) ([]int, error) {
	var skipList []int

	if err := c.Database.Table(Campaign).Select("id").Where("approvalStatus = ? OR syncTime >= ?", true, time.Now().UTC().Add(-1*time.Hour)).Find(&skipList).Error; err != nil {
		log.Errorf("fail to get campaign skip sync list: %s", err.Error())
		return skipList, customError.New(customError.DatabaseError)
	}

	return skipList, nil
}
