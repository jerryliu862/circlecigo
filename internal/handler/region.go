package handler

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"17live_wso_be/internal/service"
	"17live_wso_be/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListRegion(c *gin.Context) {
	regionsClaim, _ := c.Get("regions")
	regions, _ := regionsClaim.([]string)

	rp, err := service.New().ListRegion(c, regions)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, rp)
}

func ListTaxRate(c *gin.Context) {
	rp, err := service.New().ListTaxRate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, rp)
}

func SetTaxRate(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	var data []model.RegionWithTaxList
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Warnf("invalid request data: %v. %s", data, err.Error())
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	for _, region := range data {
		if err := util.ValidateData(region); err != nil {
			log.Warnf("invalid request data: %v. %s", data, err.Error())
			c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
			return
		}

		if _, err := service.New().GetRegion(c, region.Code); err != nil {
			log.Warnf("region does not exist: %s", region.Code)
			c.JSON(http.StatusNotFound, customError.New(customError.InvalidRequestData))
			return
		}
	}

	if err := service.New().SetTaxRate(c, data, claim.Id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, EmptyResp{})
}
