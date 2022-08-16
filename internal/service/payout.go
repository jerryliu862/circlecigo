package service

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"context"
	"time"
)

func (c *Client) ListPayout(ctx context.Context, filter model.PayoutFilter) ([]model.PayoutList, int, error) {
	log.Infof("service list payout with filter: %v", filter)
	return c.RepositoryClient.ListPayout(ctx, filter)
}

func (c *Client) ListPayoutReport(ctx context.Context, filter model.PayoutFilter) ([]model.PayoutReport, int, error) {
	log.Infof("service list payout report with filter: %v", filter)
	return c.RepositoryClient.ListPayoutReport(ctx, filter)
}

func (c *Client) GetPayoutById(ctx context.Context, id int) ([]model.Payout, error) {
	log.Infof("service get payout by id: %d", id)
	return c.RepositoryClient.GetPayoutById(ctx, id)
}

func (c *Client) GetPayoutOutline(ctx context.Context, id int) ([]model.PayoutList, error) {
	log.Infof("service get payout outline: %d", id)
	return c.RepositoryClient.GetPayoutOutline(ctx, id)
}

func (c *Client) GetPayoutReportOutline(ctx context.Context, id int) ([]model.PayoutReport, error) {
	log.Infof("service get payout report outline: %d", id)
	return c.RepositoryClient.GetPayoutReportOutline(ctx, id)
}

func (c *Client) GetGroupedPayoutDetail(ctx context.Context, id int, filter model.PayoutFilter) ([]model.PayoutListDetail, int, error) {
	log.Infof("service get grouped payout detail of %d with filter: %v", id, filter)
	return c.RepositoryClient.GetGroupedPayoutDetail(ctx, id, filter)
}

func (c *Client) GetUngroupedPayoutDetail(ctx context.Context, filter model.PayoutFilter) ([]model.PayoutListDetail, int, error) {
	log.Infof("service get ungrouped payout detail with filter: %v", filter)
	return c.RepositoryClient.GetUngroupedPayoutDetail(ctx, filter)
}

func (c *Client) GetPayoutReportDetail(ctx context.Context, payoutId int, filter model.PayoutFilter) ([]model.PayoutReportDetail, int, error) {
	log.Infof("service get payout report detail of %d with filter: %v", payoutId, filter)
	return c.RepositoryClient.GetPayoutReportDetail(ctx, payoutId, filter)
}

func (c *Client) GetPayoutReportDetailByAgency(ctx context.Context, payoutId int, filter model.PayoutFilter) ([]model.PayoutReportDetailByAgency, int, error) {
	log.Infof("service get payout report detail by agency of %d with filter: %v", payoutId, filter)
	return c.RepositoryClient.GetPayoutReportDetailByAgency(ctx, payoutId, filter)
}

func (c *Client) GetPayoutReportExcelData(ctx context.Context, payoutId int, regions []string) (model.PayoutReportExcel, error) {
	log.Infof("service get payout report excel data: payoutId %d, regions %v", payoutId, regions)

	var data model.PayoutReportExcel

	local, err := c.RepositoryClient.GetPayoutReportDetailByPayTypeForExcel(ctx, payoutId, regions, model.PayTypeLocal)
	if err != nil {
		return data, err
	}

	foreign, err := c.RepositoryClient.GetPayoutReportDetailByPayTypeForExcel(ctx, payoutId, regions, model.PayTypeForeign)
	if err != nil {
		return data, err
	}

	agency, err := c.RepositoryClient.GetPayoutReportDetailByAgencyForExcel(ctx, payoutId, regions)
	if err != nil {
		return data, err
	}

	transportation, err := c.RepositoryClient.GetPayoutReportDetailTransportationForExcel(ctx, payoutId, regions)
	if err != nil {
		return data, err
	}

	missing, err := c.RepositoryClient.GetPayoutReportDetailMissingForExcel(ctx, payoutId, regions)
	if err != nil {
		return data, err
	}

	data = model.PayoutReportExcel{
		SheetLocalData:          parseReportExcelDetail(ctx, local),
		SheetForeignData:        parseReportExcelDetail(ctx, foreign),
		SheetAgencyData:         parseReportExcelAgency(ctx, agency),
		SheetTransportationData: parseReportExcelDetail(ctx, transportation),
		SheetMissingData:        parseReportExcelDetail(ctx, missing),
	}

	data.GenerateExcelTop()

	return data, nil
}

func (c *Client) GetUngroupedBonusIdList(ctx context.Context, region, payDate string) ([]int, error) {
	log.Infof("service get ungrouped bonus id list: region %s, payDate %s", region, payDate)
	return c.RepositoryClient.GetUngroupedBonusIdList(ctx, region, payDate)
}

// including create payout records
func (c *Client) UpdatePayoutStatus(ctx context.Context, data model.Payout) error {
	log.Infof("service update payout status: %v", data)
	return c.RepositoryClient.UpdatePayoutStatus(ctx, data)
}

func (c *Client) GroupPayout(ctx context.Context, data model.PayoutGroup, uid int) error {
	log.Infof("service group payout: %v", data)

	if _, err := c.GetRegion(ctx, data.RegionCode); err != nil {
		log.Warnf("region does not exist: %s", data.RegionCode)
		return customError.New(customError.InvalidRequestData)
	}

	if err := checkPaydate(ctx, data.PayDate); err != nil {
		return err
	}

	if err := c.checkPayoutDetail(ctx, data.Ids); err != nil {
		return err
	}

	time := time.Now().UTC()

	payout := model.Payout{
		RegionCode:  data.RegionCode,
		PayDate:     data.PayDate,
		CreateTime:  time,
		CreatorUID:  uid,
		ModifyTime:  time,
		ModifierUID: uid,
	}

	payout, err := c.RepositoryClient.CreatePayout(ctx, payout)
	if err != nil {
		return err
	}

	if err := c.RepositoryClient.UpsertPayoutDetails(ctx, parsePayoutDetailList(ctx, payout.Id, data.Ids, time, uid)); err != nil {
		return err
	}

	return c.RepositoryClient.SetCampaignBonusPaydate(ctx, data.PayDate, data.Ids, uid)
}

func (c *Client) SetPayoutDate(ctx context.Context, data model.PayoutGroup, uid int) (bool, error) {
	log.Infof("service set payout date: %v", data)

	var noPayoutRemained bool

	if err := checkPaydate(ctx, data.PayDate); err != nil {
		return noPayoutRemained, err
	}

	if err := c.checkPayoutStatus(ctx, data.Ids); err != nil {
		return noPayoutRemained, err
	}

	for _, id := range data.Ids {
		bonus, err := c.RepositoryClient.GetCampaignBonusById(ctx, id)
		if err != nil {
			return noPayoutRemained, err
		} else if len(bonus) == 0 {
			log.Warnf("campaign bonus does not exist: %d", id)
			return noPayoutRemained, customError.New(customError.CampaignBonusNotExist)
		}

		originalPaydate := bonus[0].PayDate
		if originalPaydate != nil {
			originalPaydate = bonus[0].ParseDateFormat().PayDate
		}

		if originalPaydate == nil || data.PayDate != *originalPaydate {
			if err := c.RepositoryClient.SetCampaignBonusPaydate(ctx, data.PayDate, []int{id}, uid); err != nil {
				return noPayoutRemained, err
			} else if res, err := c.checkAndDeletePayoutDetail(ctx, id, originalPaydate); err != nil {
				return noPayoutRemained, err
			} else if !noPayoutRemained {
				noPayoutRemained = res
			}
		}
	}

	return noPayoutRemained, nil
}

func (c *Client) AdjustPayout(ctx context.Context, data model.CampaignBonus, uid int) (bool, error) {
	log.Infof("service adjust payout: %v", data)

	var deletePayout bool

	time := time.Now().UTC()

	if data.PayDate == nil {
		log.Warnf("invalid paydate, got nil pointer")
		return deletePayout, customError.New(customError.InvalidRequestData)
	} else if err := checkPaydate(ctx, *data.PayDate); err != nil {
		return deletePayout, err
	}

	// add bonus
	if data.Id == 0 {
		data.CreateTime = time
		data.CreatorUID = uid
		data.ModifyTime = time
		data.ModifierUID = uid

		return deletePayout, c.addBonus(ctx, data)
	}

	// renew bonus
	existingBonus, err := c.RepositoryClient.GetCampaignBonusByRankIdAndBonusType(ctx, data.RankID, data.BonusType)
	if err != nil {
		return deletePayout, err
	}

	if len(existingBonus) == 0 || existingBonus[0].Id != data.Id {
		log.Warnf("campaign bonus with rank id %d and bonusType %d does not exist: %d", data.RankID, data.BonusType, data.Id)
		return deletePayout, customError.New(customError.CampaignBonusNotExist)
	}

	bonus := existingBonus[0]

	originalPaydate := bonus.PayDate
	if originalPaydate != nil {
		originalPaydate = bonus.ParseDateFormat().PayDate
	}

	bonus.PayDate = data.PayDate
	bonus.Remark = data.Remark
	bonus.ModifyTime = time
	bonus.ModifierUID = uid

	if data.BonusType != model.BonusTypeFixed && data.BonusType != model.BonusTypeVariable {
		bonus.Amount = data.Amount
	}

	return c.renewBonus(ctx, bonus, originalPaydate)
}

func (c *Client) DeletePayoutAdjustment(ctx context.Context, bonusId int, uid int) (bool, error) {
	log.Infof("service delete payout adjustment of bonus: %d, user %d", bonusId, uid)

	var noPayoutRemained bool

	bonus, err := c.RepositoryClient.GetCampaignBonusById(ctx, bonusId)
	if err != nil {
		return noPayoutRemained, err
	} else if len(bonus) == 0 {
		log.Warnf("campaign bonus does not exist: %d", bonusId)
		return noPayoutRemained, customError.New(customError.CampaignBonusNotExist)
	} else if bonus[0].BonusType == model.BonusTypeFixed || bonus[0].BonusType == model.BonusTypeVariable {
		log.Warnf("cannot delete payout adjustment of fixed or variable bonus: bonusType %d, bonusId %d", bonus[0].BonusType, bonusId)
		return noPayoutRemained, customError.New(customError.PayoutAdjustmentNonDeletable)
	}

	if err := c.checkPayoutStatus(ctx, []int{bonusId}); err != nil {
		return noPayoutRemained, err
	}

	res, err := c.checkAndDeletePayoutDetail(ctx, bonusId, bonus[0].PayDate)
	if err != nil {
		return noPayoutRemained, err
	}

	noPayoutRemained = res

	return noPayoutRemained, c.RepositoryClient.DeleteCampaignBonus(ctx, bonusId)
}

func (c *Client) checkPayoutDetail(ctx context.Context, bonusIdList []int) error {
	for _, id := range bonusIdList {
		if payoutDetail, err := c.RepositoryClient.GetPayoutDetailByBonusId(ctx, id); err != nil {
			return err
		} else if len(payoutDetail) != 0 {
			log.Warnf("campaign bonus %d has already grouped", id)
			return customError.New(customError.PayoutGrouped)
		}
	}

	return nil
}

func (c *Client) checkPayoutStatus(ctx context.Context, bonusIdList []int) error {
	for _, id := range bonusIdList {
		if payout, err := c.RepositoryClient.GetPayoutByBonusId(ctx, id); err != nil {
			return err
		} else if len(payout) != 0 {
			if payout[0].PayStatus {
				log.Warnf("campaign bonus already paid: %d", id)
				return customError.New(customError.CampaignBonusPaid)
			}
		}
	}

	return nil
}

func (c *Client) checkAndDeletePayoutDetail(ctx context.Context, bonusId int, payDate *string) (bool, error) {
	var noPayoutRemained bool

	details, err := c.RepositoryClient.GetPayoutDetailByBonusId(ctx, bonusId)
	if err != nil {
		return noPayoutRemained, err
	}

	if len(details) != 0 {
		if err := c.RepositoryClient.DeletePayoutDetail(ctx, bonusId); err != nil {
			return noPayoutRemained, err
		}

		payoutId := details[0].PayoutID

		remained, err := c.RepositoryClient.GetPayoutDetailByPayoutId(ctx, payoutId)
		if err != nil {
			return noPayoutRemained, err
		} else if len(remained) != 0 {
			return noPayoutRemained, nil
		}

		log.Infof("no payout detail remains in payout %d", payoutId)

		if err := c.RepositoryClient.DeletePayout(ctx, payoutId); err != nil {
			return noPayoutRemained, err
		}

		noPayoutRemained = true

		return noPayoutRemained, nil
	}

	// handle ungrouped bonus
	campaign, err := c.RepositoryClient.GetCampaignByBonusId(ctx, bonusId)
	if err != nil {
		return noPayoutRemained, err
	}

	var checkRegion string
	if campaign.RegionCode != nil {
		checkRegion = *campaign.RegionCode
	}

	var checkDate string
	if payDate != nil {
		checkDate = *payDate
	} else {
		checkDate = campaign.PayDate
	}

	count, err := c.RepositoryClient.CountUngroupedPayoutDetailWithoutSelectedBonus(ctx, bonusId, checkRegion, checkDate)
	if err != nil {
		return noPayoutRemained, err
	} else if count == 0 {
		noPayoutRemained = true
	}

	return noPayoutRemained, nil
}

func (c *Client) addBonus(ctx context.Context, data model.CampaignBonus) error {
	fixed, err := c.RepositoryClient.GetCampaignBonusByRankIdAndBonusType(ctx, data.RankID, model.BonusTypeFixed)
	if err != nil {
		return err
	}

	if len(fixed) == 0 {
		log.Errorf("rank %d does not have fiexed bonus record", data.RankID)
		return customError.New(customError.CampaignFixedBonusNotExist)
	}

	bonus, err := c.RepositoryClient.CreateCampaignBonus(ctx, data)
	if err != nil {
		return err
	}

	if fixed[0].PayDate != nil && *fixed[0].ParseDateFormat().PayDate == *data.PayDate {
		payout, err := c.RepositoryClient.GetPayoutByBonusId(ctx, fixed[0].Id)
		if err != nil {
			return err
		}

		if len(payout) != 0 {
			pd := model.PayoutDetail{
				PayoutID:   payout[0].Id,
				BonusID:    bonus.Id,
				CreateTime: data.CreateTime,
				CreatorUID: data.CreatorUID,
			}

			if err := c.RepositoryClient.CreatePayoutDetail(ctx, pd); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) renewBonus(ctx context.Context, data model.CampaignBonus, originalPaydate *string) (bool, error) {
	var noPayoutRemained bool

	if err := c.checkPayoutStatus(ctx, []int{data.Id}); err != nil {
		return noPayoutRemained, err
	}

	if originalPaydate == nil || *data.PayDate != *originalPaydate {
		res, err := c.checkAndDeletePayoutDetail(ctx, data.Id, originalPaydate)
		if err != nil {
			return noPayoutRemained, err
		}
		noPayoutRemained = res
	}

	return noPayoutRemained, c.RepositoryClient.UpdateCampaignBonus(ctx, data)
}

func parsePayoutDetailList(ctx context.Context, payoutID int, bonusIdList []int, createTime time.Time, creator int) []model.PayoutDetail {
	var data []model.PayoutDetail

	for _, id := range bonusIdList {
		d := model.PayoutDetail{
			PayoutID:   payoutID,
			BonusID:    id,
			CreateTime: createTime,
			CreatorUID: creator,
		}

		data = append(data, d)
	}

	return data
}

func checkPaydate(ctx context.Context, paydate string) error {
	date, err := time.Parse("2006-01-02", paydate)
	if err != nil {
		log.Warnf("invalid paydate: %s", paydate)
		return customError.New(customError.InvalidRequestData)
	}

	t := time.Now().UTC()
	firstDay := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	if date.Before(firstDay) {
		log.Warnf("paydate set early than the first day of this month: paydate %s", paydate)
		return customError.New(customError.InvalidRequestData)
	}

	return nil
}

func parseReportExcelAgency(ctx context.Context, rawData []model.PayoutReportDetailByAgencyForExcel) []map[string]interface{} {
	data := make([]map[string]interface{}, 0)

	for _, agency := range rawData {
		data = append(data, map[string]interface{}{
			"agencyName":  agency.AgencyName,
			"agencyID":    agency.AgencyId,
			"amount":      agency.TotalAmount,
			"payAmount":   agency.AfterTaxAmount,
			"payoutType":  parsePayType(model.PayTypeAgency),
			"bankName":    agency.BankName,
			"bankCode":    agency.BankCode,
			"branchName":  agency.BranchName,
			"branchCode":  agency.BranchCode,
			"bankAccount": agency.BankAccount,
			"payeeName":   agency.BankAccountName,
			"feeOwner":    parseFeeOwner(agency.FeeOwnerType),
			"taxID":       agency.TaxID,
		})
	}

	return data
}

func parseReportExcelDetail(ctx context.Context, rawData []model.PayoutReportDetailForExcel) []map[string]interface{} {
	data := make([]map[string]interface{}, 0)

	for _, detail := range rawData {
		data = append(data, map[string]interface{}{
			"userID":      detail.UserID,
			"openID":      detail.OpenID,
			"amount":      detail.PayAmount,
			"taxAmount":   detail.TaxAmount,
			"payAmount":   detail.AfterTaxAmount,
			"payoutType":  parsePayType(detail.PayType),
			"bankName":    detail.BankName,
			"bankCode":    detail.BankCode,
			"branchName":  detail.BranchName,
			"branchCode":  detail.BranchCode,
			"bankAccount": detail.BankAccount,
			"payeeName":   detail.BankAccountName,
			"feeOwner":    parseFeeOwner(detail.FeeOwnerType),
			"taxID":       detail.TaxID,
		})
	}

	return data
}

func parsePayType(payType string) string {
	var res string

	switch payType {
	case model.PayTypeLocal:
		res = "Streamer-Local"
	case model.PayTypeForeign:
		res = "Streamer-Foreign"
	case model.PayTypeAgency:
		res = "Agency-Local"
	default:
		res = "Missing"
	}

	return res
}

func parseFeeOwner(feeOwner int) string {
	var res string

	switch feeOwner {
	case model.OwnerTypeStreamer:
		res = "17Live"
	case model.OwnerTypeAgency:
		res = "User"
	default:
		res = ""
	}

	return res
}
