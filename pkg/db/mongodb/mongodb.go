package mongodb

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/wonwooseo/panawa-api/pkg/db/model"
)

type Repository struct {
	logger zerolog.Logger
	cli    *mongo.Client

	database string
}

func NewRepository(baseLogger zerolog.Logger) *Repository {
	logger := baseLogger.With().Str("caller", "db/mongodb").Logger()

	url := viper.GetString("mongodb.url")
	database := viper.GetString("mongodb.database")

	cli, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create MongoDB client")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatal().Err(err).Msg("MongoDB server is not responding")
	}

	return &Repository{
		logger:   logger,
		cli:      cli,
		database: database,
	}
}

func (r *Repository) GetLatestPrice(ctx context.Context, item string) (*model.Price, error) {
	coll := r.cli.Database(r.database).Collection("date_prices")

	var p model.Price
	opts := options.FindOne().SetSort(bson.D{{"date_unix", -1}})
	if err := coll.FindOne(ctx, bson.D{{"item_code", item}}, opts).Decode(&p); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (r *Repository) GetLatestPrices(ctx context.Context, item string, size int64) ([]*model.Price, error) {
	coll := r.cli.Database(r.database).Collection("date_prices")

	opts := options.Find().SetSort(bson.D{{"date_unix", -1}}).SetLimit(size)
	cur, err := coll.Find(ctx, bson.D{{"item_code", item}}, opts)
	if err != nil {
		return nil, err
	}
	var ps []*model.Price
	for cur.Next(ctx) {
		var p model.Price
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		ps = append(ps, &p)
	}

	return ps, nil
}

func (r *Repository) GetDateRegionalMarketPrices(ctx context.Context, item, region string, dateUnix int64) ([]*model.Price, error) {
	coll := r.cli.Database(r.database).Collection("regional_market_prices")

	opts := options.Find().SetSort(bson.D{{"market_code", 1}})
	cur, err := coll.Find(ctx, bson.D{
		{"item_code", item},
		{"date_unix", dateUnix},
		{"region_code", region},
	}, opts)
	if err != nil {
		return nil, err
	}
	var ps []*model.Price
	for cur.Next(ctx) {
		var p model.Price
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		ps = append(ps, &p)
	}

	return ps, nil
}
