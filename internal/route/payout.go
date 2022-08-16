package route

import (
	"17live_wso_be/internal/handler"
	"17live_wso_be/internal/route/middleware"
	"17live_wso_be/internal/service"

	"github.com/gin-gonic/gin"
)

func Payout(r *gin.Engine, client *service.Client) {
	payout := r.Group("/payout", middleware.AuthUserToken)
	report := r.Group("/report", middleware.AuthUserToken)

	payout.GET("", middleware.ParsePayoutPageFilter, middleware.ParsePayoutFilter, middleware.AllowPayoutViewer, handler.ListPayout)
	payout.GET(":id", middleware.CheckPayoutOutline, handler.GetPayout)
	payout.GET("/:id/detail", middleware.CheckPayoutExistence, middleware.GetPayoutViewRegion, middleware.ParsePayoutPageFilter, middleware.ParsePayoutDetailFilter, middleware.AddPayoutDetailRegionFilter, middleware.AllowPayoutViewer, handler.GetGroupedPayoutDetail)
	payout.GET("/detail", middleware.ParsePayoutPageFilter, middleware.ParsePayoutDetailFilter, middleware.ParseUngroupedPayoutDetailFilter, middleware.AllowPayoutViewer, handler.GetUngroupedPayoutDetail)
	payout.POST("/group", handler.GroupPayout)
	payout.PUT("/date", handler.SetPayoutDate)
	payout.PUT("/adjustment", handler.AdjustPayout)
	payout.DELETE("/adjustment/:id", handler.DeletePayoutAdjustment)
	payout.PUT("/:id/status", middleware.CheckPayoutStatus, middleware.AllowPayoutAdmin, handler.UpdatePayoutStatus)

	report.GET("", middleware.ParsePayoutPageFilter, middleware.ParsePayoutFilter, middleware.AllowPayoutViewer, handler.ListPayoutReport)
	report.GET("/:id", middleware.CheckPayoutReportOutline, handler.GetPayoutReport)
	report.GET("/:id/detail", middleware.CheckPayoutExistence, middleware.GetPayoutViewRegion, middleware.ParsePayoutPageFilter, middleware.ParsePayoutReportDetailFilter, middleware.AddPayoutDetailRegionFilter, middleware.AllowPayoutViewer, handler.GetPayoutReportDetail)
	report.GET("/:id/excel", middleware.CheckPayoutExistence, middleware.GetPayoutViewRegion, middleware.AllowPayoutViewer, handler.DownloadPayoutReportFile)
}
