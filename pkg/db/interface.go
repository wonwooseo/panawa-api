package db

import (
	"context"

	"github.com/wonwooseo/panawa-api/pkg/db/model"
)

type Repository interface {
	GetLatestPrice(ctx context.Context, item string) (*model.Price, error)
	GetLatestPrices(ctx context.Context, item string, size int64) ([]*model.Price, error)
	GetDateRegionalMarketPrices(ctx context.Context, item, region string, dateUnix int64) ([]*model.Price, error)
}
