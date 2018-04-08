package crd

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	crdv1 "github.com/solo-io/gloo/pkg/storage/crd/solo.io/v1"
	"github.com/solo-io/gloo/pkg/protoutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func UpstreamToCrd(namespace string, upstream *v1.Upstream) (*crdv1.Upstream, error) {
	name := upstream.Name
	var status *v1.Status
	var ok bool
	if upstream.Status != nil {
		status, ok = proto.Clone(upstream.Status).(*v1.Status)
		if !ok {
			return nil, errors.New("internal error: output of proto.Clone was not expected type")
		}
	}
	var resourceVersion string
	var annotations map[string]string
	if upstream.Metadata != nil {
		resourceVersion = upstream.Metadata.ResourceVersion
		if upstream.Metadata.Namespace != "" {
			namespace = upstream.Metadata.Namespace
		}
		annotations = upstream.Metadata.Annotations
	}

	// clone and remove fields
	upstreamClone, ok := proto.Clone(upstream).(*v1.Upstream)
	if !ok {
		return nil, errors.New("internal error: output of proto.Clone was not expected type")
	}
	upstreamClone.Metadata = nil
	upstreamClone.Name = ""
	upstreamClone.Status = nil

	spec, err := protoutil.MarshalMap(upstreamClone)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert proto upstream to map[string]interface{}")
	}
	copySpec := crdv1.Spec(spec)

	return &crdv1.Upstream{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			ResourceVersion: resourceVersion,
			Annotations:     annotations,
		},
		Status: status,
		Spec:   &copySpec,
	}, nil
}

func UpstreamFromCrd(upstreamCrd *crdv1.Upstream) (*v1.Upstream, error) {
	var upstream v1.Upstream
	if upstreamCrd.Spec != nil {
		err := protoutil.UnmarshalMap(*upstreamCrd.Spec, &upstream)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert crd spec to upstream")
		}
	}
	// add removed fields to the internal object
	upstream.Name = upstreamCrd.Name
	upstream.Metadata = &v1.Metadata{
		ResourceVersion: upstreamCrd.ResourceVersion,
		Namespace:       upstreamCrd.Namespace,
		Annotations:     upstreamCrd.Annotations,
	}
	upstream.Status = upstreamCrd.Status
	return &upstream, nil
}

func VirtualHostToCrd(namespace string, virtualHost *v1.VirtualHost) (*crdv1.VirtualHost, error) {
	name := virtualHost.Name
	var status *v1.Status
	var ok bool
	if virtualHost.Status != nil {
		status, ok = proto.Clone(virtualHost.Status).(*v1.Status)
		if !ok {
			return nil, errors.New("internal error: output of proto.Clone was not expected type")
		}
	}
	var resourceVersion string
	var annotations map[string]string
	if virtualHost.Metadata != nil {
		resourceVersion = virtualHost.Metadata.ResourceVersion
		if virtualHost.Metadata.Namespace != "" {
			namespace = virtualHost.Metadata.Namespace
		}
		annotations = virtualHost.Metadata.Annotations
	}

	// clone and remove fields
	vHostClone, ok := proto.Clone(virtualHost).(*v1.VirtualHost)
	if !ok {
		return nil, errors.New("internal error: output of proto.Clone was not expected type")
	}
	vHostClone.Metadata = nil
	vHostClone.Name = ""
	vHostClone.Status = nil

	spec, err := protoutil.MarshalMap(vHostClone)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert proto upstream to map[string]interface{}")
	}
	copySpec := crdv1.Spec(spec)

	return &crdv1.VirtualHost{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			ResourceVersion: resourceVersion,
			Annotations:     annotations,
		},
		Status: status,
		Spec:   &copySpec,
	}, nil
}

func VirtualHostFromCrd(vHostCrd *crdv1.VirtualHost) (*v1.VirtualHost, error) {
	var virtualHost v1.VirtualHost
	if vHostCrd.Spec != nil {
		err := protoutil.UnmarshalMap(*vHostCrd.Spec, &virtualHost)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert crd spec to virtualhost")
		}
	}
	// add removed fields to the internal object
	virtualHost.Name = vHostCrd.Name
	virtualHost.Metadata = &v1.Metadata{
		ResourceVersion: vHostCrd.ResourceVersion,
		Namespace:       vHostCrd.Namespace,
		Annotations:     vHostCrd.Annotations,
	}
	virtualHost.Status = vHostCrd.Status
	return &virtualHost, nil
}
