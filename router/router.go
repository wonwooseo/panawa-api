package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/wonwooseo/panawa-api/pkg/db"
	"github.com/wonwooseo/panawa-api/pkg/db/mongodb"
	"github.com/wonwooseo/panawa-api/router/controller"
	"github.com/wonwooseo/panawa/pkg/code"
)

func NewRouter(baseLogger zerolog.Logger, corsAllowAll bool) *gin.Engine {
	logger := baseLogger.With().Str("caller", "router").Logger()
	router := gin.Default()

	var repo db.Repository = mongodb.NewRepository(baseLogger)

	// middlewares
	corsCfg := cors.Config{
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"Content-Type"},
		MaxAge:       1 * time.Hour,
	}
	if corsAllowAll {
		logger.Warn().Msg("CORS config set to allow all origins")
		corsCfg.AllowAllOrigins = true
	} else {
		corsCfg.AllowOrigins = []string{"http://localhost"} // TBD: domain
	}
	router.Use(cors.New(corsCfg))

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
		price.GET("/trend", priceCtrl.PriceTrendEndpoint)
		price.GET("/region", priceCtrl.RegionalPriceEndpoint)
	}

	return router
}
