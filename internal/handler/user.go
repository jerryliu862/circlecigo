package handler

import (
	"17live_wso_be/config"
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"17live_wso_be/internal/service"
	"17live_wso_be/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var data model.UserToken
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Warnf("invalid request data: %v. %s", data, err.Error())
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	rp, err := service.New().Login(c, &data)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, rp)
}

func ListUser(c *gin.Context) {
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

	rp, total, err := service.New().ListUser(c, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(total))

	c.JSON(http.StatusOK, rp)
}

func GetUserDetail(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		log.Warnf("invalid request data: %s. %s", c.Param("uid"), err.Error())
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestId))
		return
	}

	u, err := service.New().GetUser(c, model.User{Id: uid})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	} else if len(u) == 0 {
		log.Warnf("user does not exist: %d", uid)
		c.JSON(http.StatusNotFound, customError.New(customError.UserNotExist))
		return
	}

	rp, err := service.New().GetUserDetail(c, u[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, rp)
}

func CreateUser(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	var data model.UserDetail
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

	if u, err := service.New().GetUser(c, model.User{Email: data.User.Email}); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	} else if len(u) != 0 {
		log.Warnf("user email has been used: %s", data.User.Email)
		c.JSON(http.StatusConflict, customError.New(customError.UserDuplicated))
		return
	} else if !util.AdmitEmailDomain(config.New().User.Domains, data.User.Email) {
		log.Warnf("invalid user email domain: %s", data.User.Email)
		c.JSON(http.StatusBadRequest, customError.New(customError.UserEmailDomainInvalid))
		return
	}

	if err := service.New().CreateUserWithAuth(c, claim.Id, data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, EmptyResp{})
}

func UpdateUser(c *gin.Context) {
	claims, _ := c.Get("claims")
	claim, _ := claims.(*service.UserTokenClaims)

	var data model.UserDetail
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

	u, err := service.New().GetUser(c, model.User{Id: data.User.Id})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	} else if len(u) == 0 {
		log.Warnf("user does not exist: %d", data.User.Id)
		c.JSON(http.StatusNotFound, customError.New(customError.UserNotExist))
		return
	} else if data.User.Email != u[0].Email {
		log.Warnf("user email incompatible: expected %s, request as %s", u[0].Email, data.User.Email)
		c.JSON(http.StatusBadRequest, customError.New(customError.InvalidRequestData))
		return
	}

	data.User = u[0]

	if err := service.New().UpdateUserWithAuth(c, claim.Id, data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, EmptyResp{})
}
