package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/wonwooseo/panawa-api/pkg/db"
	"github.com/wonwooseo/panawa-api/pkg/db/mongodb"
	"github.com/wonwooseo/panawa-api/router/controller"
	"github.com/wonwooseo/panawa/pkg/code"
)

func NewRouter(baseLogger zerolog.Logger) *gin.Engine {
	router := gin.Default()

	var repo db.Repository = mongodb.NewRepository(baseLogger)

	// TODO: middlewares?

	// code-locale resolvers
	var itemCodeResolver code.Resolver = code.NewItemCodeResolver()
	var regionCodeResolver code.Resolver = code.NewRegionCodeResolver()
	var marketCodeResolver code.Resolver = code.NewMarketCodeResolver()

	// admin endpoints
	adminCtrl := controller.NewAdminController(baseLogger)
	router.GET("/health", adminCtrl.HealthCheckEndpoint)
	router.GET("/version", adminCtrl.VersionEndpoint)

	// price endpoints
	priceCtrl := controller.NewPriceController(baseLogger, repo, itemCodeResolver, regionCodeResolver, marketCodeResolver)
	price := router.Group("/price")
	{
		price.GET("", priceCtrl.LatestPriceEndpoint)
		price.GET("/", priceCtrl.LatestPriceEndpoint)
		price.GET("/trend", priceCtrl.PriceTrendEndpoint)
		price.GET("/region", priceCtrl.RegionalPriceEndpoint)
	}

	return router
}
