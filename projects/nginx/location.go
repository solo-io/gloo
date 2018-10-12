package nginx

type Location struct {
	Prefix    string
	Root      string
	ProxyPass string
}
