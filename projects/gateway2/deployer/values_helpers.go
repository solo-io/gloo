package deployer

import (
	"os"
	"strings"

	"github.com/solo-io/gloo/projects/gateway2/ports"
	"github.com/solo-io/gloo/projects/gloo/constants"
	"golang.org/x/exp/slices"
	api "sigs.k8s.io/gateway-api/apis/v1"
)

// This file contains helper functions that generate helm values in the format needed
// by the deployer.

func getPortsValues(gw *api.Gateway) []helmPort {
	gwPorts := []helmPort{}
	for _, l := range gw.Spec.Listeners {
		listenerPort := uint16(l.Port)
		if slices.IndexFunc(gwPorts, func(p helmPort) bool { return *p.Port == listenerPort }) != -1 {
			continue
		}
		targetPort := ports.TranslatePort(listenerPort)
		portName := string(l.Name)
		protocol := "TCP"

		var port helmPort
		port.Port = &listenerPort
		port.TargetPort = &targetPort
		port.Name = &portName
		port.Protocol = &protocol
		gwPorts = append(gwPorts, port)
	}
	return gwPorts
}

func getDeployerImageValues() *helmImage {
	image := os.Getenv(constants.GlooGatewayDeployerImage)
	defaultImageValues := &helmImage{
		// If tag is not defined, we fall back to the default behavior, which is to use that Chart version
	}

	if image == "" {
		// If the env is not defined, return the default
		return defaultImageValues
	}

	imageParts := strings.Split(image, ":")
	if len(imageParts) != 2 {
		// If the user provided an invalid override, fallback to the default
		return defaultImageValues
	}
	return &helmImage{
		Repository: &imageParts[0],
		Tag:        &imageParts[1],
	}
}
