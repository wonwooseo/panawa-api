package db

import (
	"context"
	"time"

	"github.com/wonwooseo/panawa-api/pkg/db/model"
)

type Repository interface {
	GetDatePrice(ctx context.Context, date time.Time, item string) (*model.Price, error)
	GetDateRangePrices(ctx context.Context, sDate, eDate time.Time, item string) ([]*model.Price, error)
	GetRegionalPrice(ctx context.Context, date time.Time, item, region string) (*model.Price, error)
	GetRegionalMarketPrices(ctx context.Context, date time.Time, item, region string) ([]*model.Price, error)
}
