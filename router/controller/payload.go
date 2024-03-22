package controller

// admin

type Version struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
}

// price

type Price struct {
	Low     int `json:"low"`
	Average int `json:"average"`
	High    int `json:"high"`
}

type TodayPrice struct {
	Price
	LastUpdateTime string `json:"last_update_time"`
}

type PricePerDate struct {
	Price
	Date string `json:"date"`
}

type PricePerRegion struct {
	Price
	PerMarket map[string]Price `json:"per_market"`
}
