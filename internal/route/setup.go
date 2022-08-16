package route

import (
	"17live_wso_be/internal/service"

	"github.com/gin-gonic/gin"
)

func SetUp(r *gin.Engine, client *service.Client) {
	User(r, client)
	Region(r, client)
	Campaign(r, client)
	Payout(r, client)
	Sync(r, client)
}
