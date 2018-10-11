package nginx

// TODO(talnordan): Can `.Prefix` or `.Root` be optional?
const locationContextTemplateText = `
location {{.Prefix}} {
    root {{.Root}};
}
`

func GenerateLocationContext(location *Location) ([]byte, error) {
	return generateContext(locationContextTemplateText, location)
}
