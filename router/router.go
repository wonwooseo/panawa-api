package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/wonwooseo/panawa-api/pkg/db"
	"github.com/wonwooseo/panawa-api/pkg/db/mongodb"
	"github.com/wonwooseo/panawa-api/router/controller"
	"github.com/wonwooseo/panawa/pkg/code"
)

func NewRouter(baseLogger zerolog.Logger) *gin.Engine {
	logger := baseLogger.With().Str("caller", "router").Logger()
	router := gin.Default()

	var repo db.Repository = mongodb.NewRepository(baseLogger)

	// middlewares
	corsCfg := cors.Config{
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"Content-Type"},
		AllowOrigins: viper.GetStringSlice("cors"),
		MaxAge:       1 * time.Hour,
	}
	logger.Info().Any("cors_allow_origins", corsCfg.AllowOrigins).Msg("cors allow origins")
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
