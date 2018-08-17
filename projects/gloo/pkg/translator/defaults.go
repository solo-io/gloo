package translator

import "time"

const (
	ClusterConnectionTimeout = time.Second * 5

	SslCertificateChainKey = "tls.crt"
	SslPrivateKeyKey       = "tls.key"
	SslRootCaKey           = "tls.root"
)
