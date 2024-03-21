package controller

type TodayPrice struct {
	AveragePrice   int    `json:"average_price"`
	LastUpdateTime string `json:"last_update_time"`
}

type PricePerDate struct {
	Date         string `json:"date"`
	AveragePrice int    `json:"average_price"`
}

type RegionData struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type PricePerMarket struct {
	Market       string `json:"market"`
	AveragePrice int    `json:"average_price"`
}

type PricePerRegion struct {
	AveragePrice int                       `json:"average_price"`
	PerMarket    map[string]PricePerMarket `json:"per_market"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
