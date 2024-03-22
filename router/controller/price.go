package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/wonwooseo/panawa-api/code"
	"github.com/wonwooseo/panawa-api/pkg/db"
)

type PriceController struct {
	logger          zerolog.Logger
	serverTZ        *time.Location
	defaultItemCode string

	repo               db.Repository
	itemCodeResolver   code.Resolver
	regionCodeResolver code.Resolver
	marketCodeResolver code.Resolver
}

func NewPriceController(
	baseLogger zerolog.Logger,
	repo db.Repository,
	itemResolver, regionResolver, marketResolver code.Resolver,
) *PriceController {
	return &PriceController{
		logger:             baseLogger.With().Str("caller", "controller/price").Logger(),
		serverTZ:           time.FixedZone("KST", 9*60*60), // UTC+9, make configurable?
		defaultItemCode:    "0000",                         // config?
		repo:               repo,
		itemCodeResolver:   itemResolver,
		regionCodeResolver: regionResolver,
		marketCodeResolver: marketResolver,
	}
}

const (
	queryKeyItemCode = "item_code"
)

func (c *PriceController) TodayPriceEndpoint(ctx *gin.Context) {
	itemCode := ctx.DefaultQuery(queryKeyItemCode, c.defaultItemCode)
	if _, ok := c.itemCodeResolver.ResolveCode(itemCode); !ok {
		ctx.JSON(http.StatusBadRequest, Error{
			Code:    "1001",
			Message: fmt.Sprintf("unknown item code: %s", itemCode),
		})
		return
	}
	now := time.Now().UTC().In(c.serverTZ)

	price, err := c.repo.GetDatePrice(ctx, now, itemCode)
	if err != nil {
		c.logger.Error().Err(err).Str("item_code", itemCode).Msg("failed to get price")
		ctx.JSON(http.StatusInternalServerError, Error{
			Code:    "0000",
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, TodayPrice{
		Price: Price{
			Low:     price.Low,
			Average: price.Average,
			High:    price.High,
		},
		LastUpdateTime: price.StringDateFmt(time.RFC3339),
	})
}

const (
	queryKeyTrendWindow = "window"
	trendWindowWeek     = "week"
	trendWindowMonth    = "month"
)

func (c *PriceController) PriceTrendEndpoint(ctx *gin.Context) {
	now := time.Now().UTC().In(c.serverTZ)

	itemCode := ctx.DefaultQuery(queryKeyItemCode, c.defaultItemCode)
	if _, ok := c.itemCodeResolver.ResolveCode(itemCode); !ok {
		ctx.JSON(http.StatusBadRequest, Error{
			Code:    "1001",
			Message: fmt.Sprintf("unknown item code: %s", itemCode),
		})
		return
	}
	window := ctx.DefaultQuery(queryKeyTrendWindow, trendWindowWeek)
	var sTime time.Time
	switch window {
	case trendWindowWeek:
		sTime = now.AddDate(0, 0, -7)
	case trendWindowMonth:
		sTime = now.AddDate(0, -1, 0)
	default:
		ctx.JSON(http.StatusBadRequest, Error{
			Code:    "1000",
			Message: fmt.Sprintf("unsupported window: %s", window),
		})
		return
	}

	prices, err := c.repo.GetDateRangePrices(ctx, sTime, now, itemCode)
	if err != nil {
		c.logger.Error().Err(err).Str("item_code", itemCode).Msg("failed to get prices")
		ctx.JSON(http.StatusInternalServerError, Error{
			Code:    "0000",
			Message: "internal server error",
		})
		return
	}

	res := make([]PricePerDate, len(prices))
	for i, price := range prices {
		res[i] = PricePerDate{
			Price: Price{
				Low:     price.Low,
				Average: price.Average,
				High:    price.High,
			},
			Date: price.StringDateFmt("2006-01-02"),
		}
	}

	ctx.JSON(http.StatusOK, res)
}

const (
	queryKeyRegion = "region_code"
)

func (c *PriceController) RegionalPriceEndpoint(ctx *gin.Context) {
	now := time.Now().UTC().In(c.serverTZ)

	itemCode := ctx.DefaultQuery(queryKeyItemCode, c.defaultItemCode)
	if _, ok := c.itemCodeResolver.ResolveCode(itemCode); !ok {
		ctx.JSON(http.StatusBadRequest, Error{
			Code:    "1001",
			Message: fmt.Sprintf("unknown item code: %s", itemCode),
		})
		return
	}
	regionCode := ctx.Query(queryKeyRegion)
	if _, ok := c.regionCodeResolver.ResolveCode(regionCode); !ok {
		ctx.JSON(http.StatusBadRequest, Error{
			Code:    "1002",
			Message: fmt.Sprintf("unknown region code: %s", regionCode),
		})
		return
	}

	regionalPrice, err := c.repo.GetRegionalPrice(ctx, now, itemCode, regionCode)
	if err != nil {
		c.logger.Error().Err(err).Str("item_code", itemCode).Str("region_code", regionCode).Msg("failed to get regional price")
		ctx.JSON(http.StatusInternalServerError, Error{
			Code:    "0000",
			Message: "internal server error",
		})
		return
	}
	marketPrices, err := c.repo.GetRegionalMarketPrices(ctx, now, itemCode, regionCode)
	if err != nil {
		c.logger.Error().Err(err).Str("item_code", itemCode).Str("region_code", regionCode).Msg("failed to get regional market prices")
		ctx.JSON(http.StatusInternalServerError, Error{
			Code:    "0000",
			Message: "internal server error",
		})
		return
	}
	perMarket := map[string]Price{}
	for _, price := range marketPrices {
		if price.MarketCode == nil {
			c.logger.Warn().Any("price", price).Msg("queried market price has nil market code")
			continue
		}
		perMarket[*price.MarketCode] = Price{
			Low:     price.Low,
			Average: price.Average,
			High:    price.High,
		}
	}

	ctx.JSON(http.StatusOK, PricePerRegion{
		Price: Price{
			Low:     regionalPrice.Low,
			Average: regionalPrice.Average,
			High:    regionalPrice.High,
		},
		PerMarket: perMarket,
	})
}
