package upstream

import (
	"strings"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/spf13/cobra"
)

func BindFlagsToMetadata(cmd *cobra.Command, meta *core.Metadata) {
	cmd.PersistentFlags().StringVar(&meta.Name, "metadata.name", "default", "name for this resource")
	cmd.PersistentFlags().StringVar(&meta.Namespace, "metadata.namespace", "default", "namespace for this resource")
	cmd.PersistentFlags().Var(&MapStringStringValue{meta.Labels}, "metadata.namespace", "labels for this resource. used for selection")
	cmd.PersistentFlags().Var(&MapStringStringValue{meta.Annotations}, "metadata.annotations", "annotations for this resource. used opaque key-value data")
}

type MapStringStringValue struct {
	Value map[string]string
}

func (v *MapStringStringValue) String() string {
	var pairs []string
	for k, v := range v.Value {
		pairs = append(pairs, k+"="+v)
	}
	return "[" + strings.Join(pairs, ";") + "]"
}

func (v *MapStringStringValue) Set(str string) error {
	v.Value = make(map[string]string)
	str = strings.TrimPrefix(str, "[")
	str = strings.TrimSuffix(str, "]")
	pairs := strings.Split(str, ";")
	for _, pair := range pairs {
		split := strings.Split(pair, "=")
		if len(split) != 2 {
			return errors.Errorf("invalid string %v, must be format KEY=VALUE", pair)
		}
		v.Value[split[0]] = split[1]
	}
	return nil
}

func (v *MapStringStringValue) Type() string {
	return "MapStringStringValue"
}

func BindFlagsToUpstream(cmd *cobra.Command, upstream *v1.Upstream) {
	BindFlagsToMetadata(cmd, &upstream.Metadata)
	BindFlagsToUpstreamSpec(cmd, upstream.UpstreamType)
}

func BindFlagsToUpstreamSpec(cmd *cobra.Command, upstreamSpec *v1.UpstreamSpec) {
	upstreamSpecOptions := &UpstreamSpecOptions{
		UpstreamSpec:    upstreamSpec,
		UpstreamSpecAws: &v1.UpstreamSpec_Aws{},
	}
	cmd.PersistentFlags().Var(upstreamSpecOptions, "upstream_type", "Each upstream in Gloo has a type. Supported types include `static`, `kubernetes`, `aws`, `consul`, and more.Each upstream type is handled by a corresponding Gloo plugin.")
	cmd.PersistentFlags().StringVar(&upstreamSpecOptions.UpstreamSpecAws.Aws.Region, "upstream_type.aws.region", "", "The AWS Region in which to run Lambda Functions")
	cmd.PersistentFlags().StringVar(&upstreamSpecOptions.UpstreamSpecAws.Aws.SecretRef, "upstream_type.aws.region", "", "secret ref")
}

type UpstreamSpecOptions struct {
	FlagName        string
	Selected        string
	UpstreamSpec    *v1.UpstreamSpec
	UpstreamSpecAws *v1.UpstreamSpec_Aws
}

func (v *UpstreamSpecOptions) String() string {
	return v.Selected
}

func (v *UpstreamSpecOptions) Set(str string) error {
	v.Selected = str
	switch v.Selected {
	case "aws":
		v.UpstreamSpec.UpstreamType = v.UpstreamSpecAws
	}
	return errors.Errorf("invalid selection for --%v: %v", v.FlagName, str)
}

func (v *UpstreamSpecOptions) Type() string {
	return "UpstreamSpecOptions"
}
