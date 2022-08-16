package route

import (
	"17live_wso_be/internal/handler"
	"17live_wso_be/internal/route/middleware"
	"17live_wso_be/internal/service"

	"github.com/gin-gonic/gin"
)

func Campaign(r *gin.Engine, client *service.Client) {
	campaign := r.Group("/campaign", middleware.AuthUserToken)

	campaign.GET("", middleware.ParseCampaignFilter, middleware.AllowCampaignViewer, handler.ListCampaign)
	campaign.GET("/:id", middleware.CheckCampaign, middleware.AllowCampaignViewer, handler.GetCampaignDetail)
	campaign.PUT("/:id/bonus", middleware.CheckCampaign, middleware.AllowCampaignEditor, handler.SetCampaignBonus)
	campaign.PUT("/:id/approve", middleware.CheckCampaignBeforeApproveBonus, middleware.AllowCampaignApprover, handler.ApproveCampaignBonus)
	campaign.GET("/noRegion", middleware.AllowAdmin, handler.ListNoRegionCampaign)
	campaign.PUT("/region", middleware.AllowAdmin, handler.SetCampaignRegion)
	campaign.GET("/rank", handler.ListUnpaidCampaignBonus)
}
