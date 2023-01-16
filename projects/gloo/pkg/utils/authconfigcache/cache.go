package authconfigcache

import (
	"runtime"
	"unsafe"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"

	errors "github.com/rotisserie/eris"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"k8s.io/utils/lru"
)

type AuthConfigCache struct {
	cache *lru.Cache
}

func NewAuthConfigCache(size int) *AuthConfigCache {
	// we only ever expect to need to look at the most recent snapshot.
	// we use a cache of size 10 to have some wiggle room in case we need more for a short amount of time
	// a runtime finalizer is used to make sure we release memory as soon as it is not needed.
	return &AuthConfigCache{cache: lru.New(size)}
}

type ACKey struct{ Name, Namespace string }
type ACMap map[ACKey]*extauth.AuthConfig

func (c *AuthConfigCache) FindAuthConfig(snap *v1snap.ApiSnapshot, namespace, name string) (*extauth.AuthConfig, error) {
	var authConfigs ACMap
	key := uintptr(unsafe.Pointer(snap))
	if cached, ok := c.cache.Get(key); ok {
		authConfigs = cached.(ACMap)
	} else {
		authConfigs = CollectAuthConfigs(snap.AuthConfigs)
		c.cache.Add(key, authConfigs)
		runtime.SetFinalizer(snap, c.finalize)
	}
	ac, ok := authConfigs[ACKey{Name: name, Namespace: namespace}]
	if !ok {
		return nil, errors.Errorf("list did not find authConfig %v.%v", namespace, name)
	}
	return ac, nil
}

func (c *AuthConfigCache) Len() int {
	return c.cache.Len()
}

func (c *AuthConfigCache) finalize(snap *v1snap.ApiSnapshot) {
	key := uintptr(unsafe.Pointer(snap))
	c.cache.Remove(key)
}

func CollectAuthConfigs(authConfigs extauth.AuthConfigList) ACMap {
	authConfigMap := ACMap{}
	for _, ac := range authConfigs {
		authConfigMap[ACKey{Name: ac.GetMetadata().GetName(), Namespace: ac.GetMetadata().GetNamespace()}] = ac
	}
	return authConfigMap
}
