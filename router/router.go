package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/wonwooseo/panawa/pkg/code"
	kokr "github.com/wonwooseo/panawa/pkg/code/ko/kr"

	"github.com/wonwooseo/panawa-api/pkg/db"
	"github.com/wonwooseo/panawa-api/pkg/db/mock"
	"github.com/wonwooseo/panawa-api/router/controller"
)

func NewRouter(baseLogger zerolog.Logger) *gin.Engine {
	router := gin.Default()

	// TODO: impl
	var repo db.Repository = mock.NewMockRepository()

	// TODO: middlewares?

	// code-locale resolvers
	var itemCodeResolver code.Resolver = kokr.NewItemCodeResolver()
	var regionCodeResolver code.Resolver = kokr.NewRegionCodeResolver()
	var marketCodeResolver code.Resolver = kokr.NewMarketCodeResolver()

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
