package secret

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	utilrt "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Secret map[string][]byte

type SecretRepo struct {
	repo     map[string]Secret
	lock     sync.RWMutex
	factory  informers.SharedInformerFactory
	informer cache.SharedIndexInformer
}

// NewSecretRepo returns a repository for secrets that automatically
// syncronizes with K8S
// TODO(ashish) - replace with the Secret watcher in Glue so we get
//                other secret stores beside K8S for free
func NewSecretRepo(cfg *rest.Config) (*SecretRepo, error) {
	secretRepo := &SecretRepo{repo: make(map[string]Secret)}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create client for k8s")
	}
	resyncDuration := 10 * time.Minute
	//  monitors all namespace
	factory := informers.NewSharedInformerFactory(kubeClient, resyncDuration)
	informer := factory.Core().V1().Secrets().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			secretRepo.update(obj.(*corev1.Secret))
		},
		UpdateFunc: func(old, new interface{}) {
			secretRepo.update(new.(*corev1.Secret))
		},
		DeleteFunc: func(obj interface{}) {
			secretRepo.remove(obj.(*corev1.Secret))
		},
	})
	secretRepo.factory = factory
	secretRepo.informer = informer
	return secretRepo, nil
}

func (s *SecretRepo) Run(stop <-chan struct{}) {
	defer utilrt.HandleCrash()
	go s.factory.Start(stop)
	go s.informer.Run(stop)
	s.factory.WaitForCacheSync(stop)
	log.Println("Started monitoring secrets")
}
func (s *SecretRepo) Get(key string) (Secret, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	secret, exists := s.repo[key]
	return secret, exists
}

func (s *SecretRepo) update(secret *corev1.Secret) {
	key := fmt.Sprintf("%s/%s", secret.Namespace, secret.Name)
	s.lock.Lock()
	defer s.lock.Unlock()

	s.repo[key] = secret.Data
}

func (s *SecretRepo) remove(secret *corev1.Secret) {
	key := fmt.Sprintf("%s/%s", secret.Namespace, secret.Name)
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.repo, key)
}
