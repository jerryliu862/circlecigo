package route

import (
	"17live_wso_be/internal/handler"
	"17live_wso_be/internal/route/middleware"
	"17live_wso_be/internal/service"

	"github.com/gin-gonic/gin"
)

func Sync(r *gin.Engine, client *service.Client) {
	sync := r.Group("/sync", middleware.AuthUserToken)

	sync.GET("", middleware.AllowSyncData, handler.SyncData)
}
