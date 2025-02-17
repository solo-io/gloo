package portforward

import (
	"context"
	"fmt"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"

	"k8s.io/client-go/rest"

	"github.com/avast/retry-go/v4"

	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/rotisserie/eris"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/skv2/pkg/multicluster/kubeconfig"
)

var _ PortForwarder = &apiPortForwarder{}

// NewApiPortForwarder returns an implementation of a PortForwarder that does not rely on the Kubernetes CLI
// but instead queries the Kubernetes API directly
// This implementation is preferred, but we have seen it fail occassionally with the following error:
//
//	portforward.go:394] error copying from local connection to remote stream: EOF
func NewApiPortForwarder(options ...Option) PortForwarder {
	return &apiPortForwarder{
		stopCh:     make(chan struct{}, 1),
		properties: buildPortForwardProperties(options...),

		// The following are populated when Start is invoked
		errCh:      nil,
		restConfig: nil,
	}
}

type apiPortForwarder struct {
	stopCh chan struct{}
	errCh  chan error

	// properties represents the set of user-defined values to configure the apiPortForwarder
	properties *properties

	// restConfig is the set of attributes that are passed to a Kubernetes client
	// The value is derived from the properties
	restConfig *rest.Config
}

func (f *apiPortForwarder) Start(ctx context.Context, options ...retry.Option) error {
	return retry.Do(func() error {
		return f.startOnce(ctx)
	}, options...)
}

func (f *apiPortForwarder) startOnce(ctx context.Context) error {
	logger := contextutils.LoggerFrom(ctx)

	config, err := kubeconfig.GetRestConfigWithContext(f.properties.kubeConfig, f.properties.kubeContext, "")
	if err != nil {
		return err
	}
	f.restConfig = config

	podName, err := f.getPodName(ctx)
	if err != nil {
		return err
	}

	f.errCh = make(chan error, 1)
	readyCh := make(chan struct{}, 1)

	var fw *portforward.PortForwarder
	go func() {
		for {
			select {
			case <-f.stopCh:
				return
			default:
			}
			var err error
			// Build a new port portforward.PortForwarder.
			fw, err = f.portForwarderToPod(podName, readyCh)
			if err != nil {
				f.errCh <- fmt.Errorf("building port apiPortForwarder failed: %v", err)
				return
			}
			if err = fw.ForwardPorts(); err != nil {
				f.errCh <- fmt.Errorf("port forward: %v", err)
				return
			}
			f.errCh <- nil
			// At this point, either the stopCh has been closed, or port apiPortForwarder connection is broken.
			// the port apiPortForwarder should have already been ready before.
			// No need to notify the ready channel anymore when forwarding again.
			readyCh = nil
		}
	}()

	// We want to block Start() until we have either gotten an error or have started
	// We may later get an error, but that is handled async.
	select {
	case err := <-f.errCh:
		return fmt.Errorf("failure running port forward process: %v", err)
	case <-readyCh:
		p, err := fw.GetPorts()
		if err != nil {
			return fmt.Errorf("failed to get ports: %v", err)
		}
		if len(p) == 0 {
			return fmt.Errorf("got no ports")
		}
		// Set local port now, as it may have been 0 as input
		f.properties.localPort = int(p[0].Local)
		logger.Debugf("Port forward established %v -> %v.%v:%v", f.Address(), podName, podName, f.properties.remotePort)
		// The apiPortForwarder is now ready.
		return nil
	}
}

func (f *apiPortForwarder) Address() string {
	return net.JoinHostPort(f.properties.localAddress, strconv.Itoa(f.properties.localPort))
}

// LocalPort returns the local port that is being forwarded to the remote port
func (f *apiPortForwarder) LocalPort() int {
	return f.properties.localPort
}

func (f *apiPortForwarder) Close() {
	close(f.stopCh)
	// Closing the stop channel should close anything
	// opened by f.apiPortForwarder.ForwardPorts()
}

func (f *apiPortForwarder) ErrChan() <-chan error {
	return f.errCh
}

func (f *apiPortForwarder) WaitForStop() {
	<-f.stopCh
}

func (f *apiPortForwarder) portForwarderToPod(podName string, readyCh chan struct{}) (*portforward.PortForwarder, error) {
	// the following code is based on this reference, https://github.com/kubernetes/client-go/issues/51
	roundTripper, upgrader, err := spdy.RoundTripperFor(f.restConfig)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", f.properties.resourceNamespace, podName)
	hostIP := strings.TrimLeft(f.restConfig.Host, "https:/")
	serverURL := url.URL{Scheme: "https", Path: path, Host: hostIP}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)

	return portforward.NewOnAddresses(dialer,
		[]string{f.properties.localAddress},
		[]string{fmt.Sprintf("%d:%d", f.properties.localPort, f.properties.remotePort)},
		f.stopCh,
		readyCh,
		f.properties.stdout,
		f.properties.stderr)
}

func (f *apiPortForwarder) getPodName(ctx context.Context) (string, error) {
	switch f.properties.resourceType {
	case "deploy":
		fallthrough
	case "deployment":
		pods, err := kubeutils.GetPodsForDeployment(ctx, f.restConfig, f.properties.resourceName, f.properties.resourceNamespace)
		if err != nil {
			return "", err
		}

		if len(pods) == 0 {
			return "", eris.Errorf("No pods found for deployment %s: %s", f.properties.resourceNamespace, f.properties.resourceName)
		}
		return pods[0], nil

	case "svc":
		fallthrough
	case "service":
		pods, err := kubeutils.GetPodsForService(ctx, f.restConfig, f.properties.resourceName, f.properties.resourceNamespace)
		if err != nil {
			return "", err
		}

		if len(pods) == 0 {
			return "", eris.Errorf("No pods found for service %s: %s", f.properties.resourceNamespace, f.properties.resourceName)
		}
		return pods[0], nil

	case "pod":
		return f.properties.resourceName, nil
	}

	return "", eris.Errorf("Could not determine pod name for resourceType: %s", f.properties.resourceType)
}
