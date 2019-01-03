package v1

import (
	"fmt"
	"strconv"

	"github.com/gogo/protobuf/proto"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
	gatewayV1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	glooV1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	sqoopV1 "github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"
)

type UnsupportedResourceTypeError struct {
	Kind string
}

func (e UnsupportedResourceTypeError) Error() string {
	return fmt.Sprintf("Changesets do not support resources of type [%v]", e.Kind)
}

type ChangesetResourceClientFactory struct {
	ChangesetName   string
	ChangesetClient ChangeSetClient
}

func (f *ChangesetResourceClientFactory) NewResourceClient(params factory.NewResourceClientParams) (clients.ResourceClient, error) {
	if err := validateResourceType(params.ResourceType); err != nil {
		return nil, err
	}
	return NewResourceClient(f.ChangesetClient, params.ResourceType, f.ChangesetName), nil
}

type ResourceClient struct {
	changesetName string
	csClient      ChangeSetClient
	resourceType  resources.Resource
}

func NewResourceClient(client ChangeSetClient, resourceType resources.Resource, changesetName string) *ResourceClient {
	return &ResourceClient{
		changesetName: changesetName,
		csClient:      client,
		resourceType:  resourceType,
	}
}

var _ clients.ResourceClient = &ResourceClient{}

func (rc *ResourceClient) Kind() string {
	return resources.Kind(rc.resourceType)
}

func (rc *ResourceClient) NewResource() resources.Resource {
	return resources.Clone(rc.resourceType)
}

func (rc *ResourceClient) Register() error {
	return nil
}

func (rc *ResourceClient) Read(namespace, name string, opts clients.ReadOpts) (resources.Resource, error) {

	// Retrieve the changeset
	changeset, err := rc.csClient.Read(defaults.GlooSystem, rc.changesetName, clients.ReadOpts{Ctx: opts.Ctx})
	if err != nil {
		return nil, err
	}

	// Search for the given resource in the changeset Data field
	resourceList, err := changeset.getResourceListByKind(rc.resourceType)
	if err != nil {
		return nil, err
	}
	for _, changesetResource := range resourceList {
		if changesetResource.GetMetadata().Namespace == namespace && changesetResource.GetMetadata().Name == name {
			return changesetResource, nil
		}
	}
	return nil, errors.NewNotExistErr(namespace, name)
}

func (rc *ResourceClient) Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error) {
	if err := validateResource(resource); err != nil {
		return nil, err
	}

	// Retrieve the changeset
	changeset, err := rc.csClient.Read(defaults.GlooSystem, rc.changesetName, clients.ReadOpts{Ctx: opts.Ctx})
	if err != nil {
		return nil, err
	}
	inputMeta := resource.GetMetadata()

	// If the resource already exists, verify if the overwrite operation is valid
	resourceList, err := changeset.getResourceListByKind(rc.resourceType)
	if err != nil {
		return nil, err
	}
	if changesetResource, err := resourceList.Find(inputMeta.Namespace, inputMeta.Name); err == nil {
		if !opts.OverwriteExisting {
			return nil, errors.NewExistErr(changesetResource.GetMetadata())
		}
		// Check whether version of input resource matches the one in the changeset
		if inputMeta.ResourceVersion != changesetResource.GetMetadata().ResourceVersion {
			return nil, errors.NewResourceVersionErr(inputMeta.Namespace, inputMeta.Name, inputMeta.ResourceVersion, changesetResource.GetMetadata().ResourceVersion)
		}
	}

	// Create a clone with updated metadata
	clone := proto.Clone(resource).(resources.Resource)
	inputMeta.ResourceVersion = newOrIncrementResourceVer(inputMeta.ResourceVersion)
	clone.SetMetadata(inputMeta)

	// Write the updated resource the changeset
	changeset.updateDataField(clone)

	// Increase the changeset edit count by 1
	changeset.EditCount.Value = changeset.EditCount.Value + 1

	// Update the changeset
	_, err = rc.csClient.Write(changeset, clients.WriteOpts{OverwriteExisting: true, Ctx: opts.Ctx})
	if err != nil {
		return nil, err
	}

	return clone, nil
}

func (rc *ResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	// Retrieve the changeset
	changeset, err := rc.csClient.Read(defaults.GlooSystem, rc.changesetName, clients.ReadOpts{Ctx: opts.Ctx})
	if err != nil {
		return err
	}

	resourceList, err := changeset.getResourceListByKind(rc.resourceType)
	if err != nil {
		return err
	}

	// If the resource does not exist
	if changesetResource, err := resourceList.Find(namespace, name); err != nil {
		if !opts.IgnoreNotExist {
			return errors.NewNotExistErr(namespace, name, err)
		} else {
			return nil
		}
	} else {

		// Remove resource from changeset
		changeset.removeDataField(changesetResource)

		// Increase the changeset edit count by 1
		changeset.EditCount.Value = changeset.EditCount.Value + 1

		// Update the changeset
		_, err = rc.csClient.Write(changeset, clients.WriteOpts{OverwriteExisting: true, Ctx: opts.Ctx})
		if err != nil {
			return err
		}

		return nil
	}
}

func (rc *ResourceClient) List(namespace string, opts clients.ListOpts) (resources.ResourceList, error) {

	// Retrieve the changeset
	changeset, err := rc.csClient.Read(defaults.GlooSystem, rc.changesetName, clients.ReadOpts{Ctx: opts.Ctx})
	if err != nil {
		return nil, err
	}
	return changeset.getResourceListByKind(rc.resourceType)
}

func (rc *ResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan resources.ResourceList, <-chan error, error) {

	// Start watching all changesets (currently we cannot watch only a specific one)
	changeSetChan, changeSetErrs, err := rc.csClient.Watch(defaults.GlooSystem, opts)
	if err != nil {
		return nil, nil, err
	}

	resourceListChan := make(chan resources.ResourceList)
	errorChan := make(chan error)

	go func(ns string) {
		for {
			select {
			case changesetList := <-changeSetChan:

				// If the changeset that this client is configured for is present in the list
				if cs, err := changesetList.Find(defaults.GlooSystem, rc.changesetName); err == nil {

					// Search the changeset for resources that match the given type
					resourceList, err := cs.getResourceListByKind(rc.resourceType)
					if err != nil {
						errorChan <- err
						break
					}

					// Filter by namespace
					res := resourceList.FilterByNamespaces([]string{ns})

					// Send to channel only if the list is not empty
					if len(res) > 0 {
						resourceListChan <- res
					}
				}
			case err := <-changeSetErrs:
				errorChan <- err
			case <-opts.Ctx.Done():
				close(resourceListChan)
				close(errorChan)
				return
			}
		}
	}(namespace)

	return resourceListChan, errorChan, nil
}

//noinspection GoReceiverNames
func (chg *ChangeSet) getResourceListByKind(kind resources.Resource) (resources.ResourceList, error) {
	var resourceList resources.ResourceList
	switch kind.(type) {
	case *gatewayV1.Gateway:
		for _, gateway := range chg.Data.Gateways {
			resourceList = append(resourceList, gateway)
		}
	case *gatewayV1.VirtualService:
		for _, vs := range chg.Data.VirtualServices {
			resourceList = append(resourceList, vs)
		}
	case *glooV1.Proxy:
		for _, proxy := range chg.Data.Proxies {
			resourceList = append(resourceList, proxy)
		}
	case *glooV1.Settings:
		for _, setting := range chg.Data.Settings {
			resourceList = append(resourceList, setting)
		}
	case *glooV1.Upstream:
		for _, upstream := range chg.Data.Upstreams {
			resourceList = append(resourceList, upstream)
		}
	case *sqoopV1.ResolverMap:
		for _, rm := range chg.Data.ResolverMaps {
			resourceList = append(resourceList, rm)
		}
	case *sqoopV1.Schema:
		for _, schema := range chg.Data.Schemas {
			resourceList = append(resourceList, schema)
		}
	default:
		// should never happen since we validate the resource type beforehand
		return nil, UnsupportedResourceTypeError{Kind: resources.Kind(kind)}
	}
	return resourceList, nil
}

//noinspection GoReceiverNames
func (chg *ChangeSet) updateDataField(resource resources.Resource) (resources.Resource, error) {
	inputName, inputNamespace := resource.GetMetadata().Name, resource.GetMetadata().Namespace
	switch inputRes := resource.(type) {

	case *gatewayV1.Gateway:
		for i, gateway := range chg.Data.Gateways {
			if inputName == gateway.Metadata.Name && inputNamespace == gateway.Metadata.Namespace {
				chg.Data.Gateways[i] = inputRes
				return inputRes, nil
			}
		}
		chg.Data.Gateways = append(chg.Data.Gateways, inputRes)
		return inputRes, nil

	case *gatewayV1.VirtualService:
		for i, vs := range chg.Data.VirtualServices {
			if inputName == vs.Metadata.Name && inputNamespace == vs.Metadata.Namespace {
				chg.Data.VirtualServices[i] = inputRes
				return inputRes, nil
			}
		}
		chg.Data.VirtualServices = append(chg.Data.VirtualServices, inputRes)
		return inputRes, nil

	case *glooV1.Proxy:
		for i, proxy := range chg.Data.Proxies {
			if inputName == proxy.Metadata.Name && inputNamespace == proxy.Metadata.Namespace {
				chg.Data.Proxies[i] = inputRes
				return inputRes, nil
			}
		}
		chg.Data.Proxies = append(chg.Data.Proxies, inputRes)
		return inputRes, nil

	case *glooV1.Settings:
		for i, setting := range chg.Data.Settings {
			if inputName == setting.Metadata.Name && inputNamespace == setting.Metadata.Namespace {
				chg.Data.Settings[i] = inputRes
				return inputRes, nil
			}
		}
		chg.Data.Settings = append(chg.Data.Settings, inputRes)
		return inputRes, nil

	case *glooV1.Upstream:
		for i, upstream := range chg.Data.Upstreams {
			if inputName == upstream.Metadata.Name && inputNamespace == upstream.Metadata.Namespace {
				chg.Data.Upstreams[i] = inputRes
				return inputRes, nil
			}
		}
		chg.Data.Upstreams = append(chg.Data.Upstreams, inputRes)
		return inputRes, nil

	case *sqoopV1.ResolverMap:
		for i, rm := range chg.Data.ResolverMaps {
			if inputName == rm.Metadata.Name && inputNamespace == rm.Metadata.Namespace {
				chg.Data.ResolverMaps[i] = inputRes
				return inputRes, nil
			}
		}
		chg.Data.ResolverMaps = append(chg.Data.ResolverMaps, inputRes)
		return inputRes, nil

	case *sqoopV1.Schema:
		for i, schema := range chg.Data.Schemas {
			if inputName == schema.Metadata.Name && inputNamespace == schema.Metadata.Namespace {
				chg.Data.Schemas[i] = inputRes
				return inputRes, nil
			}
		}
		chg.Data.Schemas = append(chg.Data.Schemas, inputRes)
		return inputRes, nil

	default:
		// should never happen since we validate the resource type beforehand
		return nil, UnsupportedResourceTypeError{Kind: resources.Kind(resource)}
	}
}

//noinspection GoReceiverNames
func (chg *ChangeSet) removeDataField(resource resources.Resource) error {
	inputName, inputNamespace := resource.GetMetadata().Name, resource.GetMetadata().Namespace
	switch resource.(type) {
	case *gatewayV1.Gateway:
		for i, gateway := range chg.Data.Gateways {
			if inputName == gateway.Metadata.Name && inputNamespace == gateway.Metadata.Namespace {
				chg.Data.Gateways = append(chg.Data.Gateways[:i], chg.Data.Gateways[i+1:]...)
			}
		}

	case *gatewayV1.VirtualService:
		for i, vs := range chg.Data.VirtualServices {
			if inputName == vs.Metadata.Name && inputNamespace == vs.Metadata.Namespace {
				chg.Data.VirtualServices = append(chg.Data.VirtualServices[:i], chg.Data.VirtualServices[i+1:]...)
			}
		}

	case *glooV1.Proxy:
		for i, proxy := range chg.Data.Proxies {
			if inputName == proxy.Metadata.Name && inputNamespace == proxy.Metadata.Namespace {
				chg.Data.Proxies = append(chg.Data.Proxies[:i], chg.Data.Proxies[i+1:]...)
			}
		}

	case *glooV1.Settings:
		for i, setting := range chg.Data.Settings {
			if inputName == setting.Metadata.Name && inputNamespace == setting.Metadata.Namespace {
				chg.Data.Settings = append(chg.Data.Settings[:i], chg.Data.Settings[i+1:]...)
			}
		}

	case *glooV1.Upstream:
		for i, upstream := range chg.Data.Upstreams {
			if inputName == upstream.Metadata.Name && inputNamespace == upstream.Metadata.Namespace {
				chg.Data.Upstreams = append(chg.Data.Upstreams[:i], chg.Data.Upstreams[i+1:]...)
			}
		}

	case *sqoopV1.ResolverMap:
		for i, rm := range chg.Data.ResolverMaps {
			if inputName == rm.Metadata.Name && inputNamespace == rm.Metadata.Namespace {
				chg.Data.ResolverMaps = append(chg.Data.ResolverMaps[:i], chg.Data.ResolverMaps[i+1:]...)
			}
		}

	case *sqoopV1.Schema:
		for i, schema := range chg.Data.Schemas {
			if inputName == schema.Metadata.Name && inputNamespace == schema.Metadata.Namespace {
				chg.Data.Schemas = append(chg.Data.Schemas[:i], chg.Data.Schemas[i+1:]...)
				return nil
			}
		}

	default:
		// should never happen since we validate the resource type beforehand
		return UnsupportedResourceTypeError{Kind: resources.Kind(resource)}
	}

	return nil
}

func validateResource(resource resources.Resource) error {

	// We need to validate again in case the client passes a resource of an unsupported type to the write method
	err := validateResourceType(resource)
	if err != nil {
		return err
	}

	err = resources.Validate(resource)
	if err != nil {
		return err
	}

	return nil
}

// Check whether the given resource type is supported by this client
func validateResourceType(kind resources.Resource) error {
	switch kind.(type) {
	case *gatewayV1.Gateway:
		return nil
	case *gatewayV1.VirtualService:
		return nil
	case *glooV1.Proxy:
		return nil
	case *glooV1.Settings:
		return nil
	case *glooV1.Upstream:
		return nil
	case *sqoopV1.ResolverMap:
		return nil
	case *sqoopV1.Schema:
		return nil
	default:
		return UnsupportedResourceTypeError{Kind: resources.Kind(kind)}
	}
}

func newOrIncrementResourceVer(resourceVersion string) string {
	curr, err := strconv.Atoi(resourceVersion)
	if err != nil {
		curr = 1
	}
	return fmt.Sprintf("%v", curr+1)
}
