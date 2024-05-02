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

type LatestPrice struct {
	Price
	DateUnix       int64  `json:"date_unix"`
	LastUpdateTime string `json:"last_update_time"`
}

type PricePerDate struct {
	Price
	Date string `json:"date"`
}
