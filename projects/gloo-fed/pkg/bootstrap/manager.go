package bootstrap

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	enterprisev1 "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	ratelimitv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	fedenterprisev1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1"
	fedgatewayv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1"
	fedgloov1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1"
	fedratelimitv1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1"
	v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var fedSchemes = runtime.SchemeBuilder{
	v1.AddToScheme,
	fedgloov1.AddToScheme,
	fedgatewayv1.AddToScheme,
	fedenterprisev1.AddToScheme,
	fedratelimitv1alpha1.AddToScheme,
}

var singleClusterSchemes = runtime.SchemeBuilder{
	gloov1.AddToScheme,
	gatewayv1.AddToScheme,
	enterprisev1.AddToScheme,
	ratelimitv1alpha1.AddToScheme,
	scheme.AddToScheme,
}

func MustSingleClusterManagerWithConfig(ctx context.Context, cfg *rest.Config) manager.Manager {
	die := func(err error) {
		contextutils.LoggerFrom(ctx).Fatalw("A fatal error occurred while getting single cluster manager", zap.Error(err))
	}

	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		die(err)
	}

	if err := singleClusterSchemes.AddToScheme(mgr.GetScheme()); err != nil {
		die(err)
	}

	return mgr
}

func MustLocalManagerWithConfig(ctx context.Context, cfg *rest.Config) manager.Manager {
	die := func(err error) {
		contextutils.LoggerFrom(ctx).Fatalw("A fatal error occurred while getting local manager", zap.Error(err))
	}

	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		die(err)
	}

	if err := fedSchemes.AddToScheme(mgr.GetScheme()); err != nil {
		die(err)
	}

	return mgr
}

func MustLocalManager(ctx context.Context) manager.Manager {
	die := func(err error) {
		contextutils.LoggerFrom(ctx).Fatalw("A fatal error occurred while getting local manager", zap.Error(err))
	}

	cfg, err := config.GetConfig()
	if err != nil {
		die(err)
	}

	return MustLocalManagerWithConfig(ctx, cfg)
}

func MustRemoteScheme(ctx context.Context) *runtime.Scheme {
	die := func(err error) {
		contextutils.LoggerFrom(ctx).Fatalw("A fatal error occurred while getting remote cluster scheme", zap.Error(err))
	}

	newScheme := runtime.NewScheme()
	err := gloov1.AddToScheme(newScheme)
	if err != nil {
		die(err)
	}
	err = gatewayv1.AddToScheme(newScheme)
	if err != nil {
		die(err)
	}
	err = ratelimitv1alpha1.AddToScheme(newScheme)
	if err != nil {
		die(err)
	}
	err = enterprisev1.AddToScheme(newScheme)
	if err != nil {
		die(err)
	}
	err = scheme.AddToScheme(newScheme)
	if err != nil {
		die(err)
	}
	return newScheme
}
