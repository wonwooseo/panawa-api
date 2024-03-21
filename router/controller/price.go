package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var logger = log.Logger.With().Str("caller", "controller").Logger()

func TodayPriceEndpoint(ctx *gin.Context) {
	// TODO

	ctx.JSON(http.StatusOK, TodayPrice{
		AveragePrice:   0,
		LastUpdateTime: "",
	})
}

const (
	queryKeyTrendWindow = "window"
	trendWindowWeek     = "week"
	trendWindowMonth    = "month"
)

func PriceTrendEndpoint(ctx *gin.Context) {
	window := ctx.DefaultQuery(queryKeyTrendWindow, trendWindowWeek)
	switch window {
	case trendWindowWeek:
	case trendWindowMonth:
	default:
		ctx.JSON(http.StatusBadRequest, Error{
			Code:    "1000",
			Message: fmt.Sprintf("unsupported window: %s", window),
		})
		return
	}

	// TODO

	ctx.JSON(http.StatusOK, []PricePerDate{})
}

func ListRegionsEndpoint(ctx *gin.Context) {
	// TODO

	ctx.JSON(http.StatusOK, []RegionData{})
}

func RegionalPriceEndpoint(ctx *gin.Context) {
	// TODO

	ctx.JSON(http.StatusOK, PricePerRegion{
		AveragePrice: 0,
		PerMarket:    map[string]PricePerMarket{},
	})
}
