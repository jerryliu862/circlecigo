package middleware

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"17live_wso_be/internal/service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckCampaign(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestId))
		c.Abort()
		return
	}

	u, err := service.New().GetCampaignById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	} else if len(u) == 0 {
		c.JSON(http.StatusNotFound, customError.New(customError.CampaignNotExist))
		c.Abort()
		return
	} else if u[0].RegionCode == nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.CampaignNoRegion))
		c.Abort()
		return
	}

	c.Set("campaignID", id)
	c.Set("regions", []string{*u[0].RegionCode})
}

func CheckCampaignBeforeApproveBonus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestId))
		c.Abort()
		return
	}

	u, err := service.New().GetCampaignById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	} else if len(u) == 0 {
		c.JSON(http.StatusNotFound, customError.New(customError.CampaignNotExist))
		c.Abort()
		return
	} else if u[0].RegionCode == nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.CampaignNoRegion))
		c.Abort()
		return
	} else if u[0].ApprovalStatus {
		c.JSON(http.StatusBadRequest, customError.New(customError.CampaignBonusApproved))
		c.Abort()
		return
	}

	c.Set("campaignID", id)
	c.Set("regions", []string{*u[0].RegionCode})
}

func ParseCampaignFilter(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	var f model.CampaignFilter

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

	// region
	var regions []string
	if c.Query("region") != "" {
		regions = strings.Split(c.Query("region"), "|")
		for i, region := range regions {
			regions[i] = strings.ToUpper(region)
		}
	} else {
		regions, err = service.New().GetUserAuthRegion(c, claim.Id, model.UserAuthTypeBonus, model.UserAuthLevelView)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}
	}
	f.Regions = regions

	// approval
	if c.Query("approval") != "" {
		approval, err := strconv.ParseBool(c.Query("approval"))
		if err != nil {
			c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
			c.Abort()
			return
		}
		f.Approval = &approval
	}

	// isZero
	if c.Query("isZero") != "" {
		isZero, err := strconv.ParseBool(c.Query("isZero"))
		if err != nil {
			c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
			c.Abort()
			return
		}
		f.IsZero = &isZero
	}

	// keyword
	f.Keyword = c.Query("keyword")

	c.Set("campaignFilter", f)
	c.Set("regions", regions)
}
