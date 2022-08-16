package repository

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *Client) GetRegion(ctx context.Context, code string) (model.Region, error) {
	var data model.Region

	if err := c.Database.Table(Region).First(&data, "code = ?", code).Error; err == gorm.ErrRecordNotFound {
		return data, customError.New(customError.RecordNotFound)
	} else if err != nil {
		log.Errorf("fail to get region: %s", err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetAllRegionCode(ctx context.Context) ([]string, error) {
	var data []string

	if err := c.Database.Table(Region).Where(model.Region{Show: true}).Select("code").Find(&data).Error; err != nil {
		log.Errorf("fail to get all region code: %s", err.Error())
		return nil, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetAllRegionCodeIncludingUnshowable(ctx context.Context) ([]string, error) {
	var data []string

	if err := c.Database.Table(Region).Select("code").Find(&data).Error; err != nil {
		log.Errorf("fail to get all region code: %s", err.Error())
		return nil, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) CreateRegion(ctx context.Context, data model.Region) error {
	if err := c.Database.Table(Region).Create(&data).Error; err != nil {
		log.Errorf("fail to create region: %v. %s", data, err.Error())
		return customError.New(customError.DatabaseError)
	}

	log.Infof("region created: %v", data)

	return nil
}

func (c *Client) ListAllRegion(ctx context.Context) ([]model.RegionDetail, error) {
	var data []model.RegionDetail

	if err := c.Database.Table(ViewRegionList).Find(&data).Error; err != nil {
		log.Errorf("fail to list all region: %s", err.Error())
		return nil, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) ListRegion(ctx context.Context, regions []string) ([]model.RegionDetail, error) {
	var data []model.RegionDetail

	if err := c.Database.Table(ViewRegionList).Where("code IN (?)", regions).Find(&data).Error; err != nil {
		log.Errorf("fail to list region of %v: %s", regions, err.Error())
		return nil, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) ListRegionDetail(ctx context.Context) ([]model.RegionDetail, error) {
	var data []model.RegionDetail

	if err := c.Database.Table(ViewRegionTaxList).Find(&data).Error; err != nil {
		log.Errorf("fail to list region detail: %s", err.Error())
		return nil, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) UpsertTaxRate(ctx context.Context, data []model.TaxRate) error {
	if len(data) != 0 {
		if err := c.Database.Table(TaxRate).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"taxFrom", "taxRate", "modifyTime", "modifierUID"}),
		}).Create(&data).Error; err != nil {
			log.Errorf("fail to upsert tax rate: %s", err.Error())
			return customError.New(customError.DatabaseError)
		}

		log.Infof("tax rate upserted: %v", data)
	}

	return nil
}
