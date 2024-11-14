package crd_categories

import (
	"path/filepath"

	"github.com/solo-io/skv2/codegen/util"
)

var (
	enterpriseCRDCategory = "solo-io"
	CommonCRDCategory     = "k8sgateway"

	enterpriseCRsManifest = filepath.Join(util.MustGetThisDir(), "testdata/manifests", "empty-enterprise-crs.yaml")
	ossCRsManifest        = filepath.Join(util.MustGetThisDir(), "testdata/manifests", "empty-oss-crs.yaml")

	installedEnterpriseCRs = []string{
		"authconfig.enterprise.gloo.solo.io/empty-authconfig",
		"ratelimitconfig.ratelimit.solo.io/empty-ratelimitconfig",
		"graphqlapi.graphql.gloo.solo.io/empty-graphqlapi",
	}
	installedOssCR = "virtualservice.gateway.solo.io/empty-virtualservice"
)
