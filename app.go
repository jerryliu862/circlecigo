package main

import (
	"17live_wso_be/config"
	"17live_wso_be/internal/route"
	"17live_wso_be/internal/route/middleware"
	"17live_wso_be/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.New().Mode)

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	service := service.New()
	route.SetUp(r, service)

	r.Run(fmt.Sprintf(":%d", config.New().Port))
}
