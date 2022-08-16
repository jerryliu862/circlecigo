package handler

import (
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"17live_wso_be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SyncData(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	user, err := service.New().GetUser(c, model.User{Id: claim.Id})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	} else if len(user) == 0 {
		log.Warnf("fail to get user with id: %d", claim.Id)
		c.JSON(http.StatusBadRequest, customError.New(customError.UnknownError))
		return
	}

	go service.New().SyncData(c, user[0].Email)

	c.JSON(http.StatusAccepted, EmptyResp{})
}
