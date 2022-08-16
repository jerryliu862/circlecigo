package route

import (
	"17live_wso_be/internal/handler"
	"17live_wso_be/internal/route/middleware"
	"17live_wso_be/internal/service"

	"github.com/gin-gonic/gin"
)

func Region(r *gin.Engine, client *service.Client) {
	user := r.Group("/region", middleware.AuthUserToken)

	user.GET("", middleware.CheckRegionPermission, handler.ListRegion)
	user.GET("/taxRate", middleware.AllowTaxViewer, handler.ListTaxRate)
	user.PUT("/taxRate", middleware.AllowTaxEditor, handler.SetTaxRate)
}
