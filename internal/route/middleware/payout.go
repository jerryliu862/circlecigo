package middleware

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"17live_wso_be/internal/service"
	"17live_wso_be/util"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ParsePayoutPageFilter(c *gin.Context) {
	var f model.PayoutFilter

	// pageSize
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
		c.Abort()
		return
	}
	f.PageSize = pageSize

	// pageNo
	pageNum, err := strconv.Atoi(c.Query("pageNo"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
		c.Abort()
		return
	}
	f.PageNum = pageNum

	c.Set("payoutFilter", f)
}

func ParsePayoutFilter(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	filterClaim, _ := c.Get("payoutFilter")
	f, _ := filterClaim.(model.PayoutFilter)

	// region
	var regions []string
	if c.Query("region") != "" {
		regions = strings.Split(c.Query("region"), "|")
		for i, region := range regions {
			regions[i] = strings.ToUpper(region)
		}
	} else {
		if rs, err := service.New().GetUserAuthRegion(c, claim.Id, model.UserAuthTypePayout, model.UserAuthLevelView); err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		} else {
			regions = rs
		}
	}
	f.Regions = regions

	// date
	if c.Query("date") != "" {
		date, err := time.Parse("2006-01-02", c.Query("date"))
		if err != nil {
			c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
			c.Abort()
			return
		}

		f.Date = &date
	}

	// isGrouped
	if c.Query("isGrouped") != "" {
		isGrouped, err := strconv.ParseBool(c.Query("isGrouped"))
		if err != nil {
			c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
			c.Abort()
			return
		}
		f.IsGrouped = &isGrouped
	}

	// isPaid
	if c.Query("isPaid") != "" {
		isPaid, err := strconv.ParseBool(c.Query("isPaid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
			c.Abort()
			return
		}
		f.IsPaid = &isPaid
	}

	c.Set("payoutFilter", f)
	c.Set("regions", regions)
}

func ParsePayoutDetailFilter(c *gin.Context) {
	filterClaim, _ := c.Get("payoutFilter")
	f, _ := filterClaim.(model.PayoutFilter)

	// bonusType
	bonusType, err := strconv.Atoi(c.Query("bonusType"))
	if err == nil {
		f.BonusType = &bonusType
	}

	// keyword
	f.Keyword = c.Query("keyword")

	c.Set("payoutFilter", f)
}

func ParsePayoutReportDetailFilter(c *gin.Context) {
	filterClaim, _ := c.Get("payoutFilter")
	f, _ := filterClaim.(model.PayoutFilter)

	// payType
	if c.Query("payType") == "" {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
		c.Abort()
		return
	}
	f.PayType = strings.ToUpper(c.Query("payType"))

	// keyword
	f.Keyword = c.Query("keyword")

	c.Set("payoutFilter", f)
}

func ParseUngroupedPayoutDetailFilter(c *gin.Context) {
	filterClaim, _ := c.Get("payoutFilter")
	f, _ := filterClaim.(model.PayoutFilter)

	// date
	date, err := time.Parse("2006-01-02", c.Query("payDate"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
		c.Abort()
		return
	}
	f.Date = &date

	// region
	region := c.Query("region")
	if region == "" {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
		c.Abort()
		return
	} else if _, err := service.New().GetRegion(c, strings.ToUpper(region)); err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
		return
	}
	f.Regions = []string{region}

	c.Set("payoutFilter", f)
	c.Set("regions", f.Regions)
}

func AddPayoutDetailRegionFilter(c *gin.Context) {
	filterClaim, _ := c.Get("payoutFilter")
	f, _ := filterClaim.(model.PayoutFilter)

	regionsClaim, _ := c.Get("regions")
	regions, _ := regionsClaim.([]string)

	f.Regions = regions

	c.Set("payoutFilter", f)
}

func CheckPayoutExistence(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestId))
		c.Abort()
		return
	}

	payout, err := service.New().GetPayoutById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	} else if len(payout) == 0 {
		c.JSON(http.StatusNotFound, customError.New(customError.PayoutNotExist))
		c.Abort()
		return
	}

	c.Set("payoutID", id)
	c.Set("payout", payout[0])
}

func GetPayoutViewRegion(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	regions, err := service.New().GetUserAuthRegion(c, claim.Id, model.UserAuthTypePayout, model.UserAuthLevelView)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	}

	c.Set("regions", regions)
}

func CheckPayoutStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestId))
		c.Abort()
		return
	}

	payout, err := service.New().GetPayoutById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	} else if len(payout) == 0 {
		c.JSON(http.StatusNotFound, customError.New(customError.PayoutNotExist))
		c.Abort()
		return
	} else if payout[0].PayStatus {
		c.JSON(http.StatusBadRequest, customError.New(customError.CampaignBonusPaid))
		c.Abort()
		return
	}

	c.Set("payoutID", id)
	c.Set("regions", []string{payout[0].RegionCode})
}

func CheckPayoutOutline(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestId))
		c.Abort()
		return
	}

	payout, err := service.New().GetPayoutOutline(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	} else if len(payout) == 0 {
		c.JSON(http.StatusNotFound, customError.New(customError.PayoutNotExist))
		c.Abort()
		return
	}

	regions, err := service.New().GetUserAuthRegion(c, claim.Id, model.UserAuthTypePayout, model.UserAuthLevelView)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	}

	if !util.ContainString(regions, payout[0].RegionCode) {
		c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
		c.Abort()
		return
	}

	c.Set("payout", payout[0])
	c.Set("regions", regions)
}

func CheckPayoutReportOutline(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestId))
		c.Abort()
		return
	}

	report, err := service.New().GetPayoutReportOutline(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	} else if len(report) == 0 {
		c.JSON(http.StatusNotFound, customError.New(customError.PayoutNotExist))
		c.Abort()
		return
	}

	regions, err := service.New().GetUserAuthRegion(c, claim.Id, model.UserAuthTypePayout, model.UserAuthLevelView)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	}

	if !util.ContainString(regions, report[0].RegionCode) {
		c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
		c.Abort()
		return
	}

	c.Set("report", report[0])
	c.Set("regions", regions)
}
