package fetcher

type Fetcher interface {
	FetchSecrets(secretRefs []string) (map[string]string, error)
}
