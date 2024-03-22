package kokr

type RegionCodeResolver struct {
	itemLocaleNames map[string]string
}

func NewRegionCodeResolver() RegionCodeResolver {
	return RegionCodeResolver{
		itemLocaleNames: map[string]string{
			"0000": "서울",
		},
	}
}

func (r RegionCodeResolver) SupportedCodes() map[string]string {
	return r.itemLocaleNames
}

func (r RegionCodeResolver) ResolveCode(c string) (string, bool) {
	localeName, ok := r.itemLocaleNames[c]
	return localeName, ok
}
