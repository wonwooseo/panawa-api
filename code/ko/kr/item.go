package kokr

type ItemCodeResolver struct {
	itemLocaleNames map[string]string
}

func NewItemCodeResolver() ItemCodeResolver {
	return ItemCodeResolver{
		itemLocaleNames: map[string]string{
			"0000": "대파",
		},
	}
}

func (r ItemCodeResolver) SupportedCodes() map[string]string {
	return r.itemLocaleNames
}

func (r ItemCodeResolver) ResolveCode(c string) (string, bool) {
	localeName, ok := r.itemLocaleNames[c]
	return localeName, ok
}
