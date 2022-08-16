package handler

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/excelizeLib"
	"17live_wso_be/internal/model"
	"17live_wso_be/internal/service"
	"17live_wso_be/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ListPayout(c *gin.Context) {
	payoutFilterClaim, _ := c.Get("payoutFilter")
	payoutFilter, _ := payoutFilterClaim.(model.PayoutFilter)

	rp, total, err := service.New().ListPayout(c, payoutFilter)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(total))

	c.JSON(http.StatusOK, rp)
}

func GetPayout(c *gin.Context) {
	payoutClaim, _ := c.Get("payout")
	data, _ := payoutClaim.(model.PayoutList)

	c.JSON(http.StatusOK, data)
}

func GetGroupedPayoutDetail(c *gin.Context) {
	payoutClaim, _ := c.Get("payoutID")
	id, _ := payoutClaim.(int)

	payoutFilterClaim, _ := c.Get("payoutFilter")
	payoutFilter, _ := payoutFilterClaim.(model.PayoutFilter)

	rp, total, err := service.New().GetGroupedPayoutDetail(c, id, payoutFilter)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(total))

	c.JSON(http.StatusOK, rp)
}

func GetUngroupedPayoutDetail(c *gin.Context) {
	payoutFilterClaim, _ := c.Get("payoutFilter")
	payoutFilter, _ := payoutFilterClaim.(model.PayoutFilter)

	rp, total, err := service.New().GetUngroupedPayoutDetail(c, payoutFilter)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(total))

	c.JSON(http.StatusOK, rp)
}

func GroupPayout(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	var data model.PayoutGroup
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Warnf("invalid request data: %v. %s", data, err.Error())
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	if len(data.IdList) != 0 {
		data = *data.ParseBonusIdList()
		if err := service.New().CheckCampaignBonusExistence(c, data.Ids); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
	} else {
		Ids, err := service.New().GetUngroupedBonusIdList(c, data.RegionCode, data.PayDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		data.Ids = Ids
	}

	if !service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit) {
		if !service.New().CheckCampaignBonusPermission(c, claim.Id, model.UserAuthTypePayout, model.UserAuthLevelEdit, data.Ids) {
			c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
			return
		}
	}

	if err := service.New().GroupPayout(c, data, claim.Id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, EmptyResp{})
}

func SetPayoutDate(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	var data model.PayoutGroup
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Warnf("invalid request data: %v. %s", data, err.Error())
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	data = *data.ParseBonusIdList()

	if err := service.New().CheckCampaignBonusExistence(c, data.Ids); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if !service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit) {
		if !service.New().CheckCampaignBonusPermission(c, claim.Id, model.UserAuthTypePayout, model.UserAuthLevelEdit, data.Ids) {
			c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
			return
		}
	}

	if noPayoutRemained, err := service.New().SetPayoutDate(c, data, claim.Id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	} else if noPayoutRemained {
		c.JSON(http.StatusNoContent, EmptyResp{})
		return
	}

	c.JSON(http.StatusOK, EmptyResp{})
}

func AdjustPayout(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	var data model.CampaignBonus
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Warnf("invalid request data: %v. %s", data, err.Error())
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	if err := util.ValidateData(data); err != nil {
		log.Warnf("invalid request data: %v. %s", data, err.Error())
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	if (data.BonusType == model.BonusTypeFixed || data.BonusType == model.BonusTypeVariable) && data.Id == 0 {
		log.Warnf("cannot add fixed or variable bonus when adjust payout: %v", data)
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	if (data.BonusType == model.BonusTypeFixed || data.BonusType == model.BonusTypeVariable) && data.Amount != 0 {
		log.Warnf("fixed or variable bonus with non-zero amount when adjust payout: %v", data)
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	if !service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit) {
		if !service.New().CheckCampaignRankPermission(c, claim.Id, model.UserAuthTypePayout, model.UserAuthLevelEdit, data.RankID) {
			log.Warnf("user %d does not have permission of campaign rank %d", claim.Id, data.RankID)
			c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
			return
		}
	}

	if noPayoutRemained, err := service.New().AdjustPayout(c, data, claim.Id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	} else if noPayoutRemained {
		c.JSON(http.StatusNoContent, EmptyResp{})
		return
	}

	c.JSON(http.StatusOK, EmptyResp{})
}

func UpdatePayoutStatus(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	payoutClaim, _ := c.Get("payoutID")
	id, _ := payoutClaim.(int)

	var data model.Payout
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Warnf("invalid request data: %v. %s", data, err.Error())
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	if !data.PayStatus {
		log.Warn("invalid to update payout status as false")
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	data.Id = id
	data.ModifyTime = time.Now().UTC()
	data.ModifierUID = claim.Id

	if err := service.New().UpdatePayoutStatus(c, data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, EmptyResp{})
}

func DeletePayoutAdjustment(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	bonusId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestId))
		return
	}

	if err := service.New().CheckCampaignBonusExistence(c, []int{bonusId}); err != nil {
		c.JSON(http.StatusNotFound, customError.New(customError.CampaignBonusNotExist))
		return
	}

	if !service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit) {
		if !service.New().CheckCampaignBonusPermission(c, claim.Id, model.UserAuthTypePayout, model.UserAuthLevelEdit, []int{bonusId}) {
			c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
			return
		}
	}

	if noPayoutRemained, err := service.New().DeletePayoutAdjustment(c, bonusId, claim.Id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	} else if noPayoutRemained {
		c.JSON(http.StatusNoContent, EmptyResp{})
		return
	}

	c.JSON(http.StatusOK, EmptyResp{})
}

func ListPayoutReport(c *gin.Context) {
	payoutFilterClaim, _ := c.Get("payoutFilter")
	payoutFilter, _ := payoutFilterClaim.(model.PayoutFilter)

	rp, total, err := service.New().ListPayoutReport(c, payoutFilter)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(total))

	c.JSON(http.StatusOK, rp)
}

func GetPayoutReport(c *gin.Context) {
	reportClaim, _ := c.Get("report")
	data, _ := reportClaim.(model.PayoutReport)

	c.JSON(http.StatusOK, data)
}

func GetPayoutReportDetail(c *gin.Context) {
	payoutClaim, _ := c.Get("payoutID")
	id, _ := payoutClaim.(int)

	payoutFilterClaim, _ := c.Get("payoutFilter")
	payoutFilter, _ := payoutFilterClaim.(model.PayoutFilter)

	if payoutFilter.PayType == model.PayTypeAgency {
		rp, total, err := service.New().GetPayoutReportDetailByAgency(c, id, payoutFilter)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.Header("X-Total-Count", strconv.Itoa(total))
		c.JSON(http.StatusOK, rp)
	} else {
		rp, total, err := service.New().GetPayoutReportDetail(c, id, payoutFilter)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.Header("X-Total-Count", strconv.Itoa(total))
		c.JSON(http.StatusOK, rp)
	}
}

func DownloadPayoutReportFile(c *gin.Context) {
	payoutClaim, _ := c.Get("payout")
	payout, _ := payoutClaim.(model.Payout)

	regionsClaim, _ := c.Get("regions")
	regions, _ := regionsClaim.([]string)

	data, err := service.New().GetPayoutReportExcelData(c, payout.Id, regions)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	data.Region = payout.RegionCode
	data.PayMonth = payout.PayDate[:7]

	excel := excelizeLib.NewReportExcel(data)
	excel.ExportToWeb(c)
}
