package repository

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *Client) ListPayout(ctx context.Context, filter model.PayoutFilter) ([]model.PayoutList, int, error) {
	data := make([]model.PayoutList, 0)
	var total int

	if err := c.Database.Raw("CALL p_GetPayoutList(?, ?, ?, ?, ?, ?)", filter.Date, strings.Join(filter.Regions, "|"), filter.IsGrouped, filter.IsPaid, filter.PageSize, filter.PageSize*(filter.PageNum-1)).Scan(&data).Error; err != nil {
		log.Errorf("fail to list payout with filter %v: %s", filter, err.Error())
		return data, total, customError.New(customError.DatabaseError)
	}

	for i := range data {
		data[i] = *data[i].ParseDateFormat()
	}

	if len(data) != 0 {
		total = data[0].TotalCount
	}

	return data, total, nil
}

func (c *Client) ListPayoutReport(ctx context.Context, filter model.PayoutFilter) ([]model.PayoutReport, int, error) {
	var data []model.PayoutReport
	var total int

	filterSQL := parsePayoutReportFilterSQL(ctx, filter)

	if err := c.Database.Raw("SELECT *, totalAmount - totalTaxAmount AS afterTaxAmount, (SELECT COUNT(*) FROM v_payoutReportList WHERE regionCode IN (?)"+filterSQL+") AS totalCount FROM v_payoutReportList WHERE regionCode IN (?)"+filterSQL+" ORDER BY payoutID LIMIT ? OFFSET ?", filter.Regions, filter.Regions, filter.PageSize, filter.PageSize*(filter.PageNum-1)).Find(&data).Error; err != nil {
		log.Errorf("fail to list payout report with filter %v: %s", filter, err.Error())
		return data, total, customError.New(customError.DatabaseError)
	}

	for i := range data {
		data[i] = *data[i].ParseDateFormat()
	}

	if len(data) != 0 {
		total = data[0].TotalCount
	}

	return data, total, nil
}

func (c *Client) GetPayoutOutline(ctx context.Context, id int) ([]model.PayoutList, error) {
	data := make([]model.PayoutList, 0)

	if err := c.Database.Raw("CALL p_GetPayout(?)", id).Scan(&data).Error; err != nil {
		log.Errorf("fail to get payout outline: %d. %s", id, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	for i := range data {
		data[i] = *data[i].ParseDateFormat()
	}

	return data, nil
}

func (c *Client) GetPayoutReportOutline(ctx context.Context, id int) ([]model.PayoutReport, error) {
	var data []model.PayoutReport

	if err := c.Database.Table(ViewPayoutReportList).Where("payoutID = ?", id).Select("*, totalAmount - totalTaxAmount AS afterTaxAmount").Find(&data).Error; err != nil {
		log.Errorf("fail to get payout report outline: %d. %s", id, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	for i := range data {
		data[i] = *data[i].ParseDateFormat()
	}

	return data, nil
}

func (c *Client) GetPayoutById(ctx context.Context, id int) ([]model.Payout, error) {
	var data []model.Payout

	if err := c.Database.Table(Payout).Where("id = ?", id).Find(&data).Error; err != nil {
		log.Errorf("fail to get payout: %d. %s", id, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetPayoutByBonusId(ctx context.Context, id int) ([]model.Payout, error) {
	var data []model.Payout

	if err := c.Database.Raw("SELECT p.* FROM Payout p LEFT JOIN PayoutDetail pd on pd.payoutID = p.id WHERE pd.bonusID = ?", id).Find(&data).Error; err != nil {
		log.Errorf("fail to get payout by bonus id: %d. %s", id, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetPayoutDetailByBonusId(ctx context.Context, id int) ([]model.PayoutDetail, error) {
	var data []model.PayoutDetail

	if err := c.Database.Table(PayoutDetail).Where("bonusID = ?", id).Find(&data).Error; err != nil {
		log.Errorf("fail to get payout detail by bonus id: %d. %s", id, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetPayoutDetailByPayoutId(ctx context.Context, id int) ([]model.PayoutDetail, error) {
	var data []model.PayoutDetail

	if err := c.Database.Table(PayoutDetail).Where("payoutID = ?", id).Find(&data).Error; err != nil {
		log.Errorf("fail to get payout detail by payout id: %d. %s", id, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetGroupedPayoutDetail(ctx context.Context, id int, filter model.PayoutFilter) ([]model.PayoutListDetail, int, error) {
	var data []model.PayoutListDetail
	var total int

	filterSQL := parsePayoutDetailFilterSQL(ctx, filter)

	if err := c.Database.Raw("SELECT *, (SELECT COUNT(*) FROM v_payoutDetail WHERE payoutID = ? AND regionCode IN (?)"+filterSQL+") AS totalCount FROM v_payoutDetail WHERE payoutID = ? AND regionCode IN (?)"+filterSQL+" LIMIT ? OFFSET ?", id, filter.Regions, id, filter.Regions, filter.PageSize, filter.PageSize*(filter.PageNum-1)).Find(&data).Error; err != nil {
		log.Errorf("fail to get grouped payout detail with filter %v: %s", filter, err.Error())
		return data, total, customError.New(customError.DatabaseError)
	}

	if len(data) != 0 {
		total = data[0].TotalCount
	}

	return data, total, nil
}

func (c *Client) GetUngroupedPayoutDetail(ctx context.Context, filter model.PayoutFilter) ([]model.PayoutListDetail, int, error) {
	var data []model.PayoutListDetail
	var total int

	filterSQL := parsePayoutDetailFilterSQL(ctx, filter)

	if err := c.Database.Raw("SELECT *, (SELECT COUNT(*) FROM v_payoutDetail2 WHERE payDate = ? AND regionCode IN (?)"+filterSQL+") AS totalCount FROM v_payoutDetail2 WHERE payDate = ? AND regionCode IN (?)"+filterSQL+" LIMIT ? OFFSET ?", filter.Date, filter.Regions, filter.Date, filter.Regions, filter.PageSize, filter.PageSize*(filter.PageNum-1)).Find(&data).Error; err != nil {
		log.Errorf("fail to get ungrouped payout detail with filter %v: %s", filter, err.Error())
		return data, total, customError.New(customError.DatabaseError)
	}

	if len(data) != 0 {
		total = data[0].TotalCount
	}

	return data, total, nil
}

func (c *Client) CountUngroupedPayoutDetailWithoutSelectedBonus(ctx context.Context, id int, region, payDate string) (int, error) {
	var count int

	if err := c.Database.Raw("SELECT COUNT(*) FROM v_payoutDetail2 WHERE bonusID <> ? AND regionCode = ? AND payDate = ?", id, region, payDate).Find(&count).Error; err != nil {
		log.Errorf("fail to count ungrouped payout detail without bonus %d, with region %s and payDate %s: %s", id, region, payDate, err.Error())
		return count, customError.New(customError.DatabaseError)
	}

	return count, nil
}

func (c *Client) GetPayoutReportDetail(ctx context.Context, id int, filter model.PayoutFilter) ([]model.PayoutReportDetail, int, error) {
	var data []model.PayoutReportDetail
	var total int

	filterSQL := parsePayoutReportDetailFilterSQL(ctx, filter)

	if err := c.Database.Raw("SELECT *, payAmount - taxAmount AS afterTaxAmount, (SELECT COUNT(*) FROM v_payoutReportDetail WHERE payoutID = ? AND regionCode IN (?)"+filterSQL+") AS totalCount FROM v_payoutReportDetail WHERE payoutID = ? AND regionCode IN (?)"+filterSQL+" ORDER BY userID LIMIT ? OFFSET ?", id, filter.Regions, id, filter.Regions, filter.PageSize, filter.PageSize*(filter.PageNum-1)).Find(&data).Error; err != nil {
		log.Errorf("fail to get payout report detail of %d with filter: %v. %s", id, filter, err.Error())
		return data, total, customError.New(customError.DatabaseError)
	}

	if len(data) != 0 {
		total = data[0].TotalCount
	}

	return data, total, nil
}

func (c *Client) GetPayoutReportDetailByAgency(ctx context.Context, id int, filter model.PayoutFilter) ([]model.PayoutReportDetailByAgency, int, error) {
	var data []model.PayoutReportDetailByAgency
	var total int

	filterStreamer := parsePayoutReportStreamerFilterSQL(ctx, filter)

	err := c.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw("SELECT *, totalAmount - totalTaxAmount AS afterTaxAmount, (SELECT COUNT(*) FROM v_payoutReportDetailByAgency WHERE payoutID = ? AND regionCode IN (?)) AS totalCount FROM v_payoutReportDetailByAgency WHERE payoutID = ? AND regionCode IN (?) ORDER BY agencyID LIMIT ? OFFSET ?", id, filter.Regions, id, filter.Regions, filter.PageSize, filter.PageSize*(filter.PageNum-1)).Find(&data).Error; err != nil {
			return err
		}

		for i := range data {
			if err := tx.Table(ViewPayoutReportDetailByAgency2).Where("agencyID = ? AND payoutID = ? AND regionCode IN (?)"+filterStreamer, data[i].AgencyId, id, filter.Regions).Find(&data[i].StreamerList).Error; err != nil {
				return err
			}

			if len(data[i].StreamerList) == 0 {
				data = append(data[:i], data[i+1:]...)
			}
		}

		return nil
	})

	if err != nil {
		log.Errorf("fail to get payout report detail by agency of %d with filter: %v. %s", id, filter, err.Error())
		return data, total, customError.New(customError.DatabaseError)
	}

	if len(data) != 0 {
		total = data[0].TotalCount
	}

	return data, total, nil
}

func (c *Client) GetPayoutReportDetailByPayTypeForExcel(ctx context.Context, id int, regions []string, payType string) ([]model.PayoutReportDetailForExcel, error) {
	var data []model.PayoutReportDetailForExcel

	if err := c.Database.Raw("SELECT *, payAmount - taxAmount AS afterTaxAmount FROM v_payoutReportDetail WHERE payoutID = ? AND regionCode IN (?) AND payType = ? ORDER BY userID", id, regions, payType).Find(&data).Error; err != nil {
		log.Errorf("fail to get payout report detail of %d with regions %v by payType: %s. %s", id, regions, payType, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetPayoutReportDetailMissingForExcel(ctx context.Context, id int, regions []string) ([]model.PayoutReportDetailForExcel, error) {
	var data []model.PayoutReportDetailForExcel

	if err := c.Database.Raw("SELECT *, payAmount - taxAmount AS afterTaxAmount FROM v_payoutReportDetail WHERE payoutID = ? AND regionCode IN (?) AND payType IS NULL ORDER BY userID", id, regions).Find(&data).Error; err != nil {
		log.Errorf("fail to get payout report detail missing of %d with regions: %v. %s", id, regions, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetPayoutReportDetailTransportationForExcel(ctx context.Context, id int, regions []string) ([]model.PayoutReportDetailForExcel, error) {
	var data []model.PayoutReportDetailForExcel

	if err := c.Database.Raw("SELECT *, payAmount - taxAmount AS afterTaxAmount FROM v_payoutReportDetailByBonusType WHERE payoutID = ? AND regionCode IN (?) AND bonusType = ? ORDER BY userID", id, regions, model.BonusTypeTransportation).Find(&data).Error; err != nil {
		log.Errorf("fail to get payout report detail transportation of %d with regions: %v. %s", id, regions, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetPayoutReportDetailByAgencyForExcel(ctx context.Context, id int, regions []string) ([]model.PayoutReportDetailByAgencyForExcel, error) {
	var data []model.PayoutReportDetailByAgencyForExcel

	if err := c.Database.Raw("SELECT *, totalAmount - totalTaxAmount AS afterTaxAmount FROM v_payoutReportDetailByAgency WHERE payoutID = ? AND regionCode IN (?) ORDER BY agencyID", id, regions).Find(&data).Error; err != nil {
		log.Errorf("fail to get payout report detail by agency of %d with regions: %v. %s", id, regions, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetUngroupedBonusIdList(ctx context.Context, region, payDate string) ([]int, error) {
	data := make([]int, 0)

	if err := c.Database.Table(ViewUngroupedPayout).Select("id").Where("regionCode = ? AND payDate = ?", region, payDate).Find(&data).Error; err != nil {
		log.Errorf("fail to get ungrouped bonus id list: %s", err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) CreatePayout(ctx context.Context, data model.Payout) (model.Payout, error) {
	if err := c.Database.Table(Payout).Create(&data).Error; err != nil {
		log.Errorf("fail to create payout: %v. %s", data, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	log.Infof("payout created: %v", data)

	return data, nil
}

// including create payout records
func (c *Client) UpdatePayoutStatus(ctx context.Context, data model.Payout) error {
	err := c.Database.Transaction(func(tx *gorm.DB) error {
		var payoutRecords []model.PayoutRecord

		if err := tx.Table(ViewPayoutRecord).Where("payoutID = ?", data.Id).Find(&payoutRecords).Error; err != nil {
			log.Errorf("fail to get payout record view of payout: %d. %s", data.Id, err.Error())
			return customError.New(customError.DatabaseError)
		}

		for i := range payoutRecords {
			payoutRecords[i].CreateTime = data.ModifyTime
			payoutRecords[i].CreatorUID = data.ModifierUID
		}

		if err := tx.Table(PayoutRecord).Create(&payoutRecords).Error; err != nil {
			return err
		}

		if err := tx.Exec("UPDATE Payout SET payStatus = ?, modifyTime = ?, modifierUID = ? WHERE id = ?", data.PayStatus, data.ModifyTime, data.ModifierUID, data.Id).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Errorf("fail to update payout status of %d: %s", data.Id, err.Error())
		return customError.New(customError.DatabaseError)
	}

	return nil
}

func (c *Client) CreatePayoutDetail(ctx context.Context, data model.PayoutDetail) error {
	if err := c.Database.Table(PayoutDetail).Create(&data).Error; err != nil {
		log.Errorf("fail to create payout detail: %v. %s", data, err.Error())
		return customError.New(customError.DatabaseError)
	}

	log.Infof("payout detail created: %v", data)

	return nil
}

func (c *Client) UpsertPayoutDetails(ctx context.Context, data []model.PayoutDetail) error {
	if len(data) != 0 {
		if err := c.Database.Table(PayoutDetail).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"createTime", "creatorUID"}),
		}).Create(&data).Error; err != nil {
			log.Errorf("fail to upsert payout details: %s", err.Error())
			return customError.New(customError.DatabaseError)
		}

		log.Infof("payout details upserted: %v", data)
	}

	return nil
}

func (c *Client) DeletePayout(ctx context.Context, id int) error {
	if err := c.Database.Table(Payout).Where("id = ?", id).Delete(&model.Payout{}).Error; err != nil {
		log.Errorf("fail to delete payout: %d. %s", id, err.Error())
		return customError.New(customError.DatabaseError)
	}

	log.Infof("payout deleted: %d", id)

	return nil
}

func (c *Client) DeletePayoutDetail(ctx context.Context, bonusId int) error {
	if err := c.Database.Table(PayoutDetail).Where("bonusID = ?", bonusId).Delete(&model.PayoutDetail{}).Error; err != nil {
		log.Errorf("fail to delete payout detail of bonus: %d. %s", bonusId, err.Error())
		return customError.New(customError.DatabaseError)
	}

	log.Infof("payout detail of bonus deleted: %d", bonusId)

	return nil
}

func parsePayoutDetailFilterSQL(ctx context.Context, filter model.PayoutFilter) string {
	var res string

	if filter.BonusType == nil && filter.Keyword == "" {
		return res
	}

	if filter.BonusType != nil {
		if *filter.BonusType == model.BonusTypeFixed {
			res += fmt.Sprintf(" AND bonusType IN (%d, %d)", model.BonusTypeFixed, model.BonusTypeVariable)
		} else {
			res += fmt.Sprintf(" AND bonusType = %d", *filter.BonusType)
		}
	}

	if filter.Keyword != "" {
		res += fmt.Sprintf(" AND (campaignTitle LIKE '%%%s%%' OR openID LIKE '%%%s%%' OR userID LIKE '%%%s%%')", filter.Keyword, filter.Keyword, filter.Keyword)
	}

	log.Infof("parse payout detail filter SQL with filter %v as: %s", filter, res)

	return res
}

func parsePayoutReportFilterSQL(ctx context.Context, filter model.PayoutFilter) string {
	var res string

	if filter.Date == nil && filter.IsPaid == nil {
		return res
	}

	if filter.Date != nil {
		res += fmt.Sprintf(" AND YEAR(payDate) = %d AND MONTH(payDate) = %d", filter.Date.Year(), int(filter.Date.Month()))
	}

	if filter.IsPaid != nil {
		res += fmt.Sprintf(" AND payStatus = %t", *filter.IsPaid)
	}

	log.Infof("parse payout report filter SQL with filter %v as: %s", filter, res)

	return res
}

func parsePayoutReportDetailFilterSQL(ctx context.Context, filter model.PayoutFilter) string {
	var res string

	if filter.PayType == model.PayTypeMissing {
		res += " AND payType IS NULL"
	} else {
		res += fmt.Sprintf(" AND payType = '%s'", filter.PayType)
	}

	if filter.Keyword == "" {
		return res
	}

	res += fmt.Sprintf(" AND (accountName LIKE '%%%s%%' OR openID LIKE '%%%s%%' OR userID LIKE '%%%s%%')", filter.Keyword, filter.Keyword, filter.Keyword)

	log.Infof("parse payout report detail filter SQL with filter %v as: %s", filter, res)

	return res
}

func parsePayoutReportStreamerFilterSQL(ctx context.Context, filter model.PayoutFilter) string {
	var res string

	if filter.Keyword == "" {
		return res
	}

	res += fmt.Sprintf(" AND (accountName LIKE '%%%s%%' OR openID LIKE '%%%s%%' OR userID LIKE '%%%s%%')", filter.Keyword, filter.Keyword, filter.Keyword)

	log.Infof("parse payout report streamer detail filter SQL with filter %v as: %s", filter, res)

	return res
}
