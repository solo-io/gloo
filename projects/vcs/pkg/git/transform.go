package git

import (
	"os"
	"path/filepath"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	gatewayV1 "github.com/solo-io/solo-projects/projects/gateway/pkg/api/v1"
	glooV1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"
	sqoopV1 "github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/constants"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/file"
)

// Copy the resources contained in the given directory into the given change set
func (r *Repository) ToChangeSetData() (*v1.Data, error) {

	namespace := defaults.GlooSystem
	data := &v1.Data{}

	// Create a file client to read from the root of the repository
	fileClient, err := file.NewFileClient(r.root)
	if err != nil {
		return data, err
	}

	// We have to check if the directory exists, otherwise List will fail
	if r.exists(constants.VirtualServiceRootDir) {
		virtualServices, err := fileClient.VirtualServiceClient.List(namespace, clients.ListOpts{})
		if err != nil {
			return data, errors.Errorf("Error reading virtual services from temp directory")
		}
		data.VirtualServices = make([]*gatewayV1.VirtualService, 0)
		for _, vs := range virtualServices {
			data.VirtualServices = append(data.VirtualServices, vs)
		}
	}

	if r.exists(constants.GatewayRootDir) {
		gateways, err := fileClient.GatewayClient.List(namespace, clients.ListOpts{})
		if err != nil {
			return data, errors.Errorf("Error reading gateways from temp directory")
		}
		data.Gateways = make([]*gatewayV1.Gateway, 0)
		for _, gateway := range gateways {
			data.Gateways = append(data.Gateways, gateway)
		}
	}

	if r.exists(constants.ProxyRootDir) {
		proxies, err := fileClient.ProxyClient.List(namespace, clients.ListOpts{})
		if err != nil {
			return data, errors.Errorf("Error reading proxies from temp directory")
		}
		data.Proxies = make([]*glooV1.Proxy, 0)
		for _, proxy := range proxies {
			data.Proxies = append(data.Proxies, proxy)
		}
	}

	if r.exists(constants.ResolverMapRootDir) {
		resolverMaps, err := fileClient.ResolverMapClient.List(namespace, clients.ListOpts{})
		if err != nil {
			return data, errors.Errorf("Error reading resolver maps from temp directory")
		}
		data.ResolverMaps = make([]*sqoopV1.ResolverMap, 0)
		for _, resolverMap := range resolverMaps {
			data.ResolverMaps = append(data.ResolverMaps, resolverMap)
		}
	}

	if r.exists(constants.SchemaRootDir) {
		schemas, err := fileClient.SchemaClient.List(namespace, clients.ListOpts{})
		if err != nil {
			return data, errors.Errorf("Error reading schemas from temp directory")
		}
		data.Schemas = make([]*sqoopV1.Schema, 0)
		for _, schema := range schemas {
			data.Schemas = append(data.Schemas, schema)
		}
	}

	if r.exists(constants.SettingsRootDir) {
		settings, err := fileClient.SettingsClient.List(namespace, clients.ListOpts{})
		if err != nil {
			return data, errors.Errorf("Error reading settings from temp directory")
		}
		data.Settings = make([]*glooV1.Settings, 0)
		for _, setting := range settings {
			data.Settings = append(data.Settings, setting)
		}
	}

	if r.exists(constants.UpstreamRootDir) {
		upstreams, err := fileClient.UpstreamClient.List(namespace, clients.ListOpts{})
		if err != nil {
			return data, errors.Errorf("Error reading settings from temp directory")
		}
		data.Upstreams = make([]*glooV1.Upstream, 0)
		for _, upstream := range upstreams {
			data.Upstreams = append(data.Upstreams, upstream)
		}
	}

	return data, nil
}

// Write the changeset data to the repository
func (r *Repository) Import(cs *v1.ChangeSet) error {

	// Create file client
	fileClient, err := file.NewFileClient(r.root)
	if err != nil {
		return err
	}

	for _, vs := range cs.Data.VirtualServices {
		_, err = fileClient.VirtualServiceClient.Write(vs, clients.WriteOpts{OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	for _, gateway := range cs.Data.Gateways {
		_, err = fileClient.GatewayClient.Write(gateway, clients.WriteOpts{OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	for _, proxy := range cs.Data.Proxies {
		_, err = fileClient.ProxyClient.Write(proxy, clients.WriteOpts{OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	for _, resolverMap := range cs.Data.ResolverMaps {
		_, err = fileClient.ResolverMapClient.Write(resolverMap, clients.WriteOpts{OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	for _, schema := range cs.Data.Schemas {
		_, err = fileClient.SchemaClient.Write(schema, clients.WriteOpts{OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	for _, setting := range cs.Data.Settings {
		_, err = fileClient.SettingsClient.Write(setting, clients.WriteOpts{OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	// TODO: do we do upstreams?
	for _, upstream := range cs.Data.Upstreams {
		_, err = fileClient.UpstreamClient.Write(upstream, clients.WriteOpts{OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) exists(dir string) bool {
	_, err := os.Stat(filepath.Join(r.root, dir, defaults.GlooSystem))
	if err != nil {
		return false
	}
	return true
}
