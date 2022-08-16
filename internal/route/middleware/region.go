package middleware

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"17live_wso_be/internal/service"
	"17live_wso_be/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CheckRegionPermission(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	var regions []string

	authTypeQuery := c.Query("authType")

	if authTypeQuery == "" {
		regions = []string{model.RegionAll}
		c.Set("regions", regions)
		return
	}

	// admin has all permission
	if service.New().PermissionCheck(c, claim.Id, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit) {
		regions = []string{model.RegionAll}
		c.Set("regions", regions)
		return
	}

	authType, err := strconv.Atoi(authTypeQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestQuery))
		c.Abort()
		return
	}

	regions, err = service.New().GetUserAuthRegion(c, claim.Id, authType, model.UserAuthLevelView)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	}

	if util.ContainString(regions, model.RegionAll) {
		regions = []string{model.RegionAll}
	}

	c.Set("regions", regions)
}
