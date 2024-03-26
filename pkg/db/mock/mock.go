package mock

import (
	"context"
	"time"

	"github.com/wonwooseo/panawa-api/pkg/db/model"
)

type MockRepository struct {
}

func NewRepository() *MockRepository {
	return &MockRepository{}
}

func (r *MockRepository) GetDatePrice(ctx context.Context, date time.Time, item string) (*model.Price, error) {
	return &model.Price{
		ItemCode:       "0000",
		DateUnix:       1711033200,
		Low:            800,
		Average:        1000,
		High:           1500,
		UpdateTimeUnix: 1711033200,
	}, nil
}

func (r *MockRepository) GetDateRangePrices(ctx context.Context, sDate, eDate time.Time, item string) ([]*model.Price, error) {
	return []*model.Price{{
		ItemCode:       "0000",
		DateUnix:       1711033200,
		Low:            800,
		Average:        1000,
		High:           1500,
		UpdateTimeUnix: 1711033200,
	}, {
		ItemCode:       "0000",
		DateUnix:       1711033200,
		Low:            850,
		Average:        1100,
		High:           1550,
		UpdateTimeUnix: 1711033200,
	}, {
		ItemCode:       "0000",
		DateUnix:       1711033200,
		Low:            830,
		Average:        1070,
		High:           1560,
		UpdateTimeUnix: 1711033200,
	}, {
		ItemCode:       "0000",
		DateUnix:       1711033200,
		Low:            880,
		Average:        1110,
		High:           1550,
		UpdateTimeUnix: 1711033200,
	}}, nil
}

func (r *MockRepository) GetRegionalMarketPrices(ctx context.Context, date time.Time, item, region string) ([]*model.Price, error) {
	market1 := "01"
	market2 := "02"
	return []*model.Price{{
		ItemCode:       "0000",
		DateUnix:       1711033200,
		Low:            800,
		Average:        1000,
		High:           1500,
		RegionCode:     &region,
		MarketCode:     &market1,
		UpdateTimeUnix: 1711033200,
	}, {
		ItemCode:       "0000",
		DateUnix:       1711033200,
		Low:            850,
		Average:        1100,
		High:           1550,
		RegionCode:     &region,
		MarketCode:     &market2,
		UpdateTimeUnix: 1711033200,
	}}, nil
}
