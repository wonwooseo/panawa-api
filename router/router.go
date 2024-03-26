package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/wonwooseo/panawa/pkg/code"
	"github.com/wonwooseo/panawa/pkg/code/ko"

	"github.com/wonwooseo/panawa-api/pkg/db"
	"github.com/wonwooseo/panawa-api/pkg/db/mongodb"
	"github.com/wonwooseo/panawa-api/router/controller"
)

func NewRouter(baseLogger zerolog.Logger) *gin.Engine {
	router := gin.Default()

	var repo db.Repository = mongodb.NewRepository(baseLogger)

	// TODO: middlewares?

	// code-locale resolvers
	var itemCodeResolver code.Resolver = ko.NewItemCodeResolver()
	var regionCodeResolver code.Resolver = ko.NewRegionCodeResolver()
	var marketCodeResolver code.Resolver = ko.NewMarketCodeResolver()

	// admin endpoints
	adminCtrl := controller.NewAdminController(baseLogger)
	router.GET("/health", adminCtrl.HealthCheckEndpoint)
	router.GET("/version", adminCtrl.VersionEndpoint)

	// price endpoints
	priceCtrl := controller.NewPriceController(baseLogger, repo, itemCodeResolver, regionCodeResolver, marketCodeResolver)
	price := router.Group("/price")
	{
		price.GET("", priceCtrl.TodayPriceEndpoint)
		price.GET("/", priceCtrl.TodayPriceEndpoint)
		price.GET("/trend", priceCtrl.PriceTrendEndpoint)
		price.GET("/region", priceCtrl.RegionalPriceEndpoint)
	}

	return router
}
