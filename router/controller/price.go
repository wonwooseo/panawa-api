package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/wonwooseo/panawa/pkg/code"

	"github.com/wonwooseo/panawa-api/pkg/db"
	rerr "github.com/wonwooseo/panawa-api/router/errors"
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
		serverTZ:           time.FixedZone("KST", 9*60*60),
		defaultItemCode:    "0000",
		repo:               repo,
		itemCodeResolver:   itemResolver,
		regionCodeResolver: regionResolver,
		marketCodeResolver: marketResolver,
	}
}

const (
	queryKeyItemCode = "item_code"
)

func (c *PriceController) LatestPriceEndpoint(ctx *gin.Context) {
	itemCode := ctx.DefaultQuery(queryKeyItemCode, c.defaultItemCode)
	if _, ok := c.itemCodeResolver.ResolveCode(itemCode); !ok {
		ctx.JSON(rerr.NewUnknownItemError(itemCode))
		return
	}

	price, err := c.repo.GetLatestPrice(ctx, itemCode)
	if err != nil {
		c.logger.Error().Err(err).Str("item_code", itemCode).Msg("failed to get price")
		ctx.JSON(rerr.NewInternalServerError())
		return
	}
	if price == nil {
		c.logger.Warn().Str("item_code", itemCode).Msg("no price data")
		ctx.JSON(rerr.NewNoPriceDataError())
		return
	}

	ctx.JSON(http.StatusOK, LatestPrice{
		Price: Price{
			Low:     price.Low,
			Average: price.Average,
			High:    price.High,
		},
		DateUnix:       price.DateUnix,
		LastUpdateTime: price.StringUpdateTimeFmt(time.RFC3339),
	})
}

const (
	queryKeyTrendWindow = "window"
	trendWindowWeek     = "week"
	trendWindowMonth    = "month"
)

func (c *PriceController) PriceTrendEndpoint(ctx *gin.Context) {
	itemCode := ctx.DefaultQuery(queryKeyItemCode, c.defaultItemCode)
	if _, ok := c.itemCodeResolver.ResolveCode(itemCode); !ok {
		ctx.JSON(rerr.NewUnknownItemError(itemCode))
		return
	}
	window := ctx.DefaultQuery(queryKeyTrendWindow, trendWindowWeek)
	var querySize int64
	switch window {
	case trendWindowWeek:
		querySize = 7
	case trendWindowMonth:
		querySize = 30
	default:
		ctx.JSON(rerr.NewInvalidQueryParamError(queryKeyTrendWindow, window))
		return
	}

	prices, err := c.repo.GetLatestPrices(ctx, itemCode, querySize)
	if err != nil {
		c.logger.Error().Err(err).Str("item_code", itemCode).Msg("failed to get prices")
		ctx.JSON(rerr.NewInternalServerError())
		return
	}

	res := make([]PricePerDate, len(prices))
	for i, price := range prices {
		res[len(prices)-1-i] = PricePerDate{
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
	queryKeyRegion   = "region_code"
	queryKeyDateUnix = "date_unix"
)

func (c *PriceController) RegionalPriceEndpoint(ctx *gin.Context) {
	itemCode := ctx.DefaultQuery(queryKeyItemCode, c.defaultItemCode)
	if _, ok := c.itemCodeResolver.ResolveCode(itemCode); !ok {
		ctx.JSON(rerr.NewUnknownItemError(itemCode))
		return
	}
	regionCode := ctx.Query(queryKeyRegion)
	if _, ok := c.regionCodeResolver.ResolveCode(regionCode); !ok {
		ctx.JSON(rerr.NewUnknownRegionError(regionCode))
		return
	}
	dateUnix, err := strconv.ParseInt(ctx.Query(queryKeyDateUnix), 10, 64)
	if err != nil {
		c.logger.Error().Err(err).Str("item_code", itemCode).Str("region_code", regionCode).Msg("failed to parse dateUnix")
		ctx.JSON(rerr.NewInvalidQueryParamError(queryKeyDateUnix, ctx.Query(queryKeyDateUnix)))
		return
	}

	marketPrices, err := c.repo.GetDateRegionalMarketPrices(ctx, itemCode, regionCode, dateUnix)
	if err != nil {
		c.logger.Error().Err(err).Str("item_code", itemCode).Str("region_code", regionCode).Msg("failed to get regional market prices")
		ctx.JSON(rerr.NewInternalServerError())
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

	ctx.JSON(http.StatusOK, perMarket)
}
