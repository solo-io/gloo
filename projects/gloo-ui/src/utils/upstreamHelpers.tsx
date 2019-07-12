import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';

/* -------------------------------------------------------------------------- */
/*                                  UPSTREAMS                                 */
/* -------------------------------------------------------------------------- */

export function getUpstreamType(upstream: Upstream.AsObject) {
  let upstreamType = '';
  if (!!upstream.upstreamSpec!.aws) {
    upstreamType = 'AWS';
  }
  if (!!upstream.upstreamSpec!.azure) {
    upstreamType = 'Azure';
  }

  if (!!upstream.upstreamSpec!.consul) {
    upstreamType = 'Consul';
  }

  if (!!upstream.upstreamSpec!.kube) {
    upstreamType = 'Kubernetes';
  }
  return upstreamType;
}

export const UPSTREAM_TYPES = [
  {
    key: 'AWS',
    value: 'AWS'
  },
  {
    key: 'Azure',
    value: 'Azure'
  },
  {
    key: 'Kubernetes',
    value: 'Kubernetes'
  },
  {
    key: 'Static',
    value: 'Static'
  },
  {
    key: 'Consul',
    value: 'Consul'
  }
];

export enum UPSTREAM_SPEC_TYPES {
  AZURE = 'Azure',
  KUBE = 'Kubernetes',
  AWS = 'AWS',
  STATIC = 'Static',
  CONSUL = 'Consul'
}

// from https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html#concepts-available-regions
export const AWS_REGIONS = [
  {
    name: 'us-east-1',
    label: 'US East(N.Virginia)'
  },
  {
    name: 'us-east-2',
    label: 'US East(Ohio)'
  },
  {
    name: 'us-west-1',
    label: 'US West(N.California)'
  },
  {
    name: 'us-west-2',
    label: 'US West(Oregon)'
  },
  {
    name: 'ca-central-1',
    label: 'Canada(Central)'
  },
  {
    name: 'eu-central-1',
    label: 'EU(Frankfurt)'
  },
  {
    name: 'eu-west-1',
    label: 'EU(Ireland)'
  },
  {
    name: 'eu-west-2',
    label: 'EU(London)'
  },
  {
    name: 'eu-west-3',
    label: 'EU(Paris)'
  },
  {
    name: 'eu-north-1',
    label: 'EU(Stockholm)'
  },
  {
    name: 'ap-east-1',
    label: 'Asia Pacific(Hong Kong)'
  },
  {
    name: 'ap-northeast-1',
    label: 'Asia Pacific(Tokyo)'
  },
  {
    name: 'ap-northeast-2',
    label: 'Asia Pacific(Seoul)'
  },
  {
    name: 'ap-northeast-3',
    label: 'Asia Pacific(Osaka-Local)'
  },
  {
    name: 'ap-southeast-1',
    label: 'Asia Pacific(Singapore)'
  },
  {
    name: 'ap-southeast-2',
    label: 'Asia Pacific(Sydney)'
  },
  {
    name: 'ap-south-1',
    label: 'Asia Pacific(Mumbai)'
  },
  {
    name: 'sa-east-1',
    label: 'South America(São Paulo)'
  }
];
