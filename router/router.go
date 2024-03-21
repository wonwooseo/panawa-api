package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonwooseo/panawa-api/router/controller"
)

func NewRouter(version, buildTime string) *gin.Engine {
	router := gin.Default()

	// healthcheck
	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	// version
	router.GET("/version", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, struct {
			Version   string
			BuildTime string
		}{
			Version:   version,
			BuildTime: buildTime,
		})
	})

	// price endpoints
	price := router.Group("/price")
	{
		// today
		price.GET("/", controller.TodayPriceEndpoint)
		// past trend
		price.GET("/trend", controller.PriceTrendEndpoint)
		// list regions
		price.GET("/region", controller.ListRegionsEndpoint)
		// per region
		price.GET("/region/:code", controller.RegionalPriceEndpoint)
	}

	return router
}
