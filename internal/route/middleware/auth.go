package middleware

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"17live_wso_be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

const Authorization = "Authorization"

func AuthUserToken(c *gin.Context) {
	token := c.Request.Header.Get(Authorization)
	claims, err := service.New().ValidUserToken(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		c.Abort()
		return
	}
	c.Set("claims", claims)
}

func AllowAdmin(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	if !service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit) {
		c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
		c.Abort()
		return
	}
}

func AllowTaxViewer(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	// admin has all permission
	if service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit) {
		return
	}

	if !service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeTax, model.UserAuthLevelView) {
		c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
		c.Abort()
		return
	}
}

func AllowTaxEditor(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	// admin has all permission
	if service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit) {
		return
	}

	if !service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeTax, model.UserAuthLevelEdit) {
		c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
		c.Abort()
		return
	}
}

func AllowSyncData(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	// admin has all permission
	if service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit) {
		return
	}

	if res, err := service.New().GetUserAuthRegion(c, claim.Id, model.UserAuthTypeBonus, model.UserAuthLevelView); err != nil || len(res) == 0 {
		c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
		c.Abort()
		return
	}
}

func AllowCampaignViewer(c *gin.Context) {
	checkRegionPermission(c, model.UserAuthTypeBonus, model.UserAuthLevelView)
}

func AllowCampaignEditor(c *gin.Context) {
	checkRegionPermission(c, model.UserAuthTypeBonus, model.UserAuthLevelEdit)
}

func AllowCampaignApprover(c *gin.Context) {
	checkRegionPermission(c, model.UserAuthTypeBonus, model.UserAuthLevelApprove)
}

func AllowPayoutViewer(c *gin.Context) {
	checkRegionPermission(c, model.UserAuthTypePayout, model.UserAuthLevelView)
}

func AllowPayoutAdmin(c *gin.Context) {
	checkRegionPermission(c, model.UserAuthTypePayout, model.UserAuthLevelEdit)
}

func checkRegionPermission(c *gin.Context, authType int, authLevel int) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	// admin has all permission
	if service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit) {
		return
	}

	regions, _ := c.Get("regions")
	regionList, _ := regions.([]string)

	if len(regionList) == 0 {
		c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
		c.Abort()
		return
	}

	for _, region := range regionList {
		if !service.New().PermissionCheck(c, claim.Id, region, authType, authLevel) {
			c.JSON(http.StatusForbidden, customError.New(customError.PermissionDenied))
			c.Abort()
			return
		}
	}
}
