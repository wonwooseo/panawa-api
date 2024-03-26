package controller

import (
	"net/http"
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
	tzYesterday := time.Now().UTC().In(c.serverTZ).AddDate(0, 0, -1)
	searchDate := time.Date(tzYesterday.Year(), tzYesterday.Month(), tzYesterday.Day(), 0, 0, 0, 0, c.serverTZ)

	itemCode := ctx.DefaultQuery(queryKeyItemCode, c.defaultItemCode)
	if _, ok := c.itemCodeResolver.ResolveCode(itemCode); !ok {
		ctx.JSON(rerr.NewUnknownItemError(itemCode))
		return
	}

	price, err := c.repo.GetDatePrice(ctx, searchDate, itemCode)
	if err != nil {
		c.logger.Error().Err(err).Str("item_code", itemCode).Msg("failed to get price")
		ctx.JSON(rerr.NewInternalServerError())
		return
	}
	if price == nil {
		c.logger.Warn().Str("item_code", itemCode).Msg("latest price is not yet fetched")
		ctx.JSON(rerr.NewNoPriceDataError())
		return
	}

	ctx.JSON(http.StatusOK, TodayPrice{
		Price: Price{
			Low:     price.Low,
			Average: price.Average,
			High:    price.High,
		},
		LastUpdateTime: price.StringUpdateTimeFmt(time.RFC3339),
	})
}

const (
	queryKeyTrendWindow = "window"
	trendWindowWeek     = "week"
	trendWindowMonth    = "month"
)

func (c *PriceController) PriceTrendEndpoint(ctx *gin.Context) {
	tzYesterday := time.Now().UTC().In(c.serverTZ).AddDate(0, 0, -1)
	eDate := time.Date(tzYesterday.Year(), tzYesterday.Month(), tzYesterday.Day(), 0, 0, 0, 0, c.serverTZ)

	itemCode := ctx.DefaultQuery(queryKeyItemCode, c.defaultItemCode)
	if _, ok := c.itemCodeResolver.ResolveCode(itemCode); !ok {
		ctx.JSON(rerr.NewUnknownItemError(itemCode))
		return
	}
	window := ctx.DefaultQuery(queryKeyTrendWindow, trendWindowWeek)
	var sDate time.Time
	switch window {
	case trendWindowWeek:
		sDate = eDate.AddDate(0, 0, -7)
	case trendWindowMonth:
		sDate = eDate.AddDate(0, -1, 0)
	default:
		ctx.JSON(rerr.NewInvalidQueryParamError(queryKeyTrendWindow, window))
		return
	}

	prices, err := c.repo.GetDateRangePrices(ctx, sDate, eDate, itemCode)
	if err != nil {
		c.logger.Error().Err(err).Str("item_code", itemCode).Msg("failed to get prices")
		ctx.JSON(rerr.NewInternalServerError())
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
	tzYesterday := time.Now().UTC().In(c.serverTZ).AddDate(0, 0, -1)
	searchDate := time.Date(tzYesterday.Year(), tzYesterday.Month(), tzYesterday.Day(), 0, 0, 0, 0, c.serverTZ)

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

	marketPrices, err := c.repo.GetRegionalMarketPrices(ctx, searchDate, itemCode, regionCode)
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
