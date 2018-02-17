package vault

import (
	"time"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/solo-io/gloo/pkg/secretwatcher"
)

type vaultSecretWatcher struct {
	secrets    chan secretwatcher.SecretMap
	errors     chan error
	secretRefs []string
	lastSeen   uint64
	client     *vaultapi.Client
}

func NewVaultSecretWatcher(resyncFrequency time.Duration, retries int, vaultAddr, token string, stopCh <-chan struct{}) (*vaultSecretWatcher, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = vaultAddr
	cfg.MaxRetries = retries
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "starting vault client")
	}
	client.SetToken(token)

	tick := time.Tick(resyncFrequency)
	w := &vaultSecretWatcher{
		secrets: make(chan secretwatcher.SecretMap),
		errors:  make(chan error),
		client:  client,
	}

	go func() {
		for {
			select {
			case <-tick:
				w.updateSecrets()
			case <-stopCh:
				return
			}
		}
	}()
	return w, nil
}

func (w *vaultSecretWatcher) updateSecrets() {
	secretMap, err := w.getSecrets()
	if err != nil {
		w.errors <- err
		return
	}
	// ignore empty configs / no secrets to watch
	if len(secretMap) == 0 {
		return
	}
	w.secrets <- secretMap
}

// triggers an update
func (w *vaultSecretWatcher) TrackSecrets(secretRefs []string) {
	w.secretRefs = secretRefs
	w.updateSecrets()
}

func (w *vaultSecretWatcher) Secrets() <-chan secretwatcher.SecretMap {
	return w.secrets
}

func (w *vaultSecretWatcher) Error() <-chan error {
	return w.errors
}
func (w *vaultSecretWatcher) getSecrets() (secretwatcher.SecretMap, error) {
	secrets := make(secretwatcher.SecretMap)
	for _, ref := range w.secretRefs {
		secret, err := w.client.Logical().Read(ref)
		if err != nil {
			return nil, errors.Wrapf(err, "reading secret %v", ref)
		}
		secretData := make(map[string]string)
		for key, value := range secret.Data {
			strValue, ok := value.(string)
			if !ok {
				return nil, errors.New("secret data must be encoded as string:string pairs")
			}
			secretData[key] = strValue
		}
		secrets[ref] = secretData
	}
	hash, err := hashstructure.Hash(secrets, nil)
	if err != nil {
		runtime.HandleError(err)
		return nil, nil
	}
	if w.lastSeen == hash {
		return nil, nil
	}
	w.lastSeen = hash
	return secrets, nil
}
