package clients

type LicensingClient interface {
	Validate(key string) (bool, error)
}
