package service

import (
	"17live_wso_be/internal/model"
	"context"
	"sort"
	"time"
)

func (c *Client) GetRegion(ctx context.Context, code string) (model.Region, error) {
	return c.RepositoryClient.GetRegion(ctx, code)
}

func (c *Client) GetAllRegionCode(ctx context.Context) ([]string, error) {
	return c.RepositoryClient.GetAllRegionCode(ctx)
}

func (c *Client) GetAllRegionCodeIncludingUnshowable(ctx context.Context) ([]string, error) {
	return c.RepositoryClient.GetAllRegionCodeIncludingUnshowable(ctx)
}

func (c *Client) CreateRegion(ctx context.Context, code string) error {
	log.Infof("service create region: %s", code)

	data := model.Region{
		Code: code,
		Name: code,
	}

	return c.RepositoryClient.CreateRegion(ctx, data)
}

func (c *Client) ListRegion(ctx context.Context, regions []string) ([]model.RegionDetail, error) {
	log.Infof("service list region: %v", regions)

	data := make([]model.RegionDetail, 0)

	if len(regions) == 0 {
		return data, nil
	}

	if regions[0] == model.RegionAll {
		return c.RepositoryClient.ListAllRegion(ctx)
	}

	return c.RepositoryClient.ListRegion(ctx, regions)
}

func (c *Client) ListTaxRate(ctx context.Context) ([]model.RegionWithTaxList, error) {
	log.Infof("service list tax rate")

	var data []model.RegionWithTaxList

	details, err := c.RepositoryClient.ListRegionDetail(ctx)
	if err != nil {
		return data, err
	}

	data = parseRegionDetail(ctx, details)

	return data, nil
}

func (c *Client) SetTaxRate(ctx context.Context, data []model.RegionWithTaxList, uid int) error {
	log.Infof("service set tax rate: %v", data)

	var taxList []model.TaxRate

	time := time.Now().UTC()

	for _, rd := range data {
		for _, t := range rd.TaxList {
			tax := model.TaxRate{
				RegionCode:  rd.Code,
				PayType:     t.PayType,
				TaxFrom:     t.TaxFrom,
				TaxRate:     t.TaxRate,
				CreateTime:  time,
				CreatorUID:  uid,
				ModifyTime:  time,
				ModifierUID: uid,
			}
			taxList = append(taxList, tax)
		}
	}

	return c.RepositoryClient.UpsertTaxRate(ctx, taxList)
}

func parseRegionDetail(ctx context.Context, details []model.RegionDetail) []model.RegionWithTaxList {
	var data []model.RegionWithTaxList

	regionMap := make(map[string]model.RegionDetail)
	taxMap := make(map[string][]model.TaxRate)

	for _, d := range details {
		if _, ok := regionMap[d.Code]; !ok {
			regionMap[d.Code] = model.RegionDetail{
				Code:           d.Code,
				Name:           d.Name,
				CurrencyCode:   d.CurrencyCode,
				CurrencyName:   d.CurrencyName,
				CurrencyFormat: d.CurrencyFormat,
			}
			taxMap[d.Code] = make([]model.TaxRate, 0)
		}

		if d.PayType != "" {
			tax := model.TaxRate{
				PayType: d.PayType,
				TaxFrom: d.TaxFrom,
				TaxRate: d.TaxRate,
			}
			taxMap[d.Code] = append(taxMap[d.Code], tax)
		}
	}

	for i := range regionMap {
		var r model.RegionWithTaxList
		r.RegionDetail = regionMap[i]
		r.TaxList = taxMap[i]
		data = append(data, r)
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Code < data[j].Code
	})

	return data
}
