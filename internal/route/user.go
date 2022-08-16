package route

import (
	"17live_wso_be/internal/handler"
	"17live_wso_be/internal/route/middleware"
	"17live_wso_be/internal/service"

	"github.com/gin-gonic/gin"
)

func User(r *gin.Engine, client *service.Client) {
	user := r.Group("/user")

	user.POST("/login", handler.Login)
	user.GET("", middleware.AuthUserToken, middleware.AllowAdmin, handler.ListUser)
	user.POST("", middleware.AuthUserToken, middleware.AllowAdmin, handler.CreateUser)
	user.PUT("", middleware.AuthUserToken, middleware.AllowAdmin, handler.UpdateUser)
	user.GET("/:uid", middleware.AuthUserToken, middleware.AllowAdmin, handler.GetUserDetail)
}
