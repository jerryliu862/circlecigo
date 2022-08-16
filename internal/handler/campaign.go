package handler

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"17live_wso_be/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListCampaign(c *gin.Context) {
	campaignFilterClaim, _ := c.Get("campaignFilter")
	campaignFilter, _ := campaignFilterClaim.(model.CampaignFilter)

	rp, total, err := service.New().ListCampaign(c, campaignFilter)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(total))

	c.JSON(http.StatusOK, rp)
}

func GetCampaignDetail(c *gin.Context) {
	campaignClaim, _ := c.Get("campaignID")
	id, _ := campaignClaim.(int)

	rp, err := service.New().GetCampaignDetail(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, rp)
}

func SetCampaignBonus(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	campaignClaim, _ := c.Get("campaignID")
	id, _ := campaignClaim.(int)

	var data model.CampaignDetail
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Warnf("invalid request data: %v. %s", data, err.Error())
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	data.Id = id

	if err := service.New().SetCampaignBonus(c, claim.Id, data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, EmptyResp{})
}

func ApproveCampaignBonus(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	campaignClaim, _ := c.Get("campaignID")
	id, _ := campaignClaim.(int)

	if err := service.New().ApproveCampaignBonus(c, id, claim.Id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, EmptyResp{})
}

func ListNoRegionCampaign(c *gin.Context) {
	var page model.PageFilter
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
		return
	}
	page.PageSize = pageSize

	pageNum, err := strconv.Atoi(c.Query("pageNo"))
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
		return
	}
	page.PageNum = pageNum

	rp, total, err := service.New().ListNoRegionCampaign(c, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(total))

	c.JSON(http.StatusOK, rp)
}

func SetCampaignRegion(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	var data []model.Campaign
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Warnf("invalid request data: %v. %s", data, err.Error())
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	if err := service.New().SetCampaignRegion(c, data, claim.Id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, EmptyResp{})
}

func ListUnpaidCampaignBonus(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	rp, err := service.New().ListUnpaidCampaignBonus(c, claim.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, rp)
}
