package nginx

// TODO(talnordan): Can `.Prefix` or `.Root` be optional?
const locationContextTemplateText = `
location {{.Prefix}} {
    root {{.Root}};
}
`

// TODO(talnordan): Deduplicate common code with `GenerateHttpContext()`
func GenerateLocationContext(location *Location) ([]byte, error) {
	return generateContext(locationContextTemplateText, location)
}
