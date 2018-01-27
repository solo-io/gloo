package types

type VirtualHost struct {
	Domains   []string
	SSLConfig SSLConfig
}

// contains certs etc.
type SSLConfig struct{}
