package kokr

type MarketCodeResolver struct {
	itemLocaleNames map[string]string
}

func NewMarketCodeResolver() MarketCodeResolver {
	return MarketCodeResolver{
		itemLocaleNames: map[string]string{
			"00": "도매시장",
			"01": "전통시장",
			"02": "대형유통",
			"03": "인터넷",
		},
	}
}

func (r MarketCodeResolver) SupportedCodes() map[string]string {
	return r.itemLocaleNames
}

func (r MarketCodeResolver) ResolveCode(c string) (string, bool) {
	localeName, ok := r.itemLocaleNames[c]
	return localeName, ok
}
