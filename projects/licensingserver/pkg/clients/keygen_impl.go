package clients

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/solo-io/keygen"
	"github.com/solo-io/keygen/services"
)

const (
	KEYGEN_CREDENTIALS_FILE = "$KEYGEN_CREDENTIALS_FILE"
	KEYGEN_ENV_EXPANSION    = "KEYGEN"
)

type KeygenAuthConfig struct {
	AccountId string `json:"accountId"envconfig:"KEYGEN_ACCOUNT_ID"`
	AuthToken string `json:"authToken"envconfig:"KEYGEN_AUTH_TOKEN"`
	Email     string `json:"email"envconfig:"KEYGEN_EMAIL"`
	Password  string `json:"password"envconfig:"KEYGEN_PASSWORD"`
}

type KeygenLicensingClient struct {
	keygenClient *keygen.LicenseClient
}

func (klc *KeygenLicensingClient) Validate(key string) (bool, error) {
	valid, err := klc.keygenClient.License.QuickValidate(key)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("&#39;%s&#39; was not found", key)) {
			return false, nil
		}
		return false, err
	}
	return valid, nil
}

func NewKeygenLicensingClient(cfg *KeygenAuthConfig) (*KeygenLicensingClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("keygen auth config cannot be nil")
	}
	if cfg.AccountId == "" || cfg.AuthToken == "" {
		return nil, fmt.Errorf("neither authToken nor accountId can be empty")
	}
	auth := services.NewTokenAuth(cfg.AuthToken)
	options := func(cli *services.WrapperClient) error {
		cli.Auth = func(r *http.Request) {
			if r.Header == nil {
				r.Header = http.Header{}
			}
			auth.AddAuthHeader(r)
		}
		cli.Account = cfg.AccountId
		return nil
	}
	lc, err := keygen.NewLicenseClient(nil, options)
	if err != nil {
		log.Fatal(err)
	}

	klc := &KeygenLicensingClient{
		keygenClient: lc,
	}
	return klc, nil
}
