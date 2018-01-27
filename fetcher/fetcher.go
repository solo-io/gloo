package fetcher

type Fetcher interface {
	FetchSecrets(names []string) (map[string]string, error)
}
