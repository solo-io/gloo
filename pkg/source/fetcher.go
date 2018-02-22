package source

// TODO(ashish)  - map to Gloo v1 api objects
type Upstream struct {
	Name      string
	Type      string
	Functions []Function
	Spec      map[string]interface{}
}

type Function struct {
	Name string
	Spec map[string]interface{}
}

// FunctionFetcher represents the function that knows how to discover
// functions for the given upstream
type FunctionFetcher interface {
	Fetch(u *Upstream) ([]Function, error)
	CanFetch(u *Upstream) bool
}

var (
	FetcherRegistry = fetcherRegistry{}
)

type fetcherRegistry struct {
	registry []FunctionFetcher
}

func (fr *fetcherRegistry) Add(f FunctionFetcher) {
	fr.registry = append(fr.registry, f)
}

func (fr *fetcherRegistry) Fetcher(u *Upstream) FunctionFetcher {
	for _, f := range fr.registry {
		if f.CanFetch(u) {
			return f
		}
	}
	return nil
}
