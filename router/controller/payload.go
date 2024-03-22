package controller

// admin

type Version struct {
	Version   string
	BuildTime string
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

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
