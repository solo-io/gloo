package types

type VirtualHost struct {
	Domain    string
	SSLConfig SSLConfig
}

// contains certs etc.
type SSLConfig struct{}
