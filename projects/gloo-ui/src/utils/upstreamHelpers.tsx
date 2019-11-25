import {
  CheckboxFilterProps,
  RadioFilterProps
} from 'Components/Common/ListingFilter';
import _ from 'lodash';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { Metadata } from 'proto/github.com/solo-io/solo-kit/api/v1/metadata_pb';
import { UpstreamInput } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
/* -------------------------------------------------------------------------- */
/*                                  UPSTREAMS                                 */
/* -------------------------------------------------------------------------- */

export function createUpstreamId(upstreamMetadata: Metadata.AsObject): string {
  return `${upstreamMetadata!.name}-.-${upstreamMetadata!.namespace}`;
}

export function parseUpstreamId(
  upstreamId: string
): {
  name: string;
  namespace: string;
} {
  const idData = upstreamId.split('-.-');

  return {
    name: idData[0],
    namespace: idData[1]
  };
}

export function getUpstreamType(upstream: Upstream.AsObject) {
  // return upstream.getUpstreamTypeCase();

  let upstreamType = 'other';
  if (upstream?.aws !== undefined) {
    upstreamType = 'Aws';
  }

  if (upstream?.azure !== undefined) {
    upstreamType = 'Azure';
  }

  if (upstream?.consul !== undefined) {
    upstreamType = 'Consul';
  }

  if (upstream?.kube !== undefined) {
    upstreamType = 'Kubernetes';
  }
  // TODO: add back
  // if (!!upstream.upstreamSpec!.awsEc2) {
  //   upstreamType = 'Aws Ec 2';
  // }

  if (upstream?.pb_static !== undefined) {
    upstreamType = 'Static';
  }

  return upstreamType;
}

export function getFunctionInfo(upstream: Upstream.AsObject) {
  if (getUpstreamType(upstream) === 'Azure') {
    return `${upstream?.azure!.functionsList.length}`;
  }
  if (getUpstreamType(upstream) === 'Aws') {
    return `${upstream?.aws!.lambdaFunctionsList.length}`;
  }
  if (upstream?.kube?.serviceSpec?.rest !== undefined) {
    return `${upstream?.kube!.serviceSpec!.rest?.transformationsMap?.length}`;
  }
  return '';
}

export function getFunctionList(upstream: Upstream.AsObject) {
  let functionsList: { key: string; value: string }[] = [];
  if (upstream) {
    if (upstream?.aws && upstream?.aws?.lambdaFunctionsList?.length > 0) {
      let newList = upstream?.aws.lambdaFunctionsList.map(lambda => {
        return {
          key: lambda.logicalName,
          value: lambda.logicalName
        };
      });
      functionsList = newList;
    }
    if (upstream?.kube) {
      const { serviceSpec } = upstream?.kube;
      if (serviceSpec && serviceSpec.rest) {
        let newList = serviceSpec.rest.transformationsMap.map(
          ([func, transform]) => {
            return {
              key: func,
              value: func
            };
          }
        );
        functionsList = newList;
      }
    }
  }
  return functionsList;
}

// The upstreams we allow the user to create are not guaranteed to be the same
// as the ones we can discover since these depend on the grpc server
export const UPSTREAM_TYPES = Object.keys(UpstreamInput.SpecCase)
  .filter(str => str !== 'SPEC_NOT_SET') // auto generated for oneof fields
  .map(str => (str === 'KUBE' ? 'Kubernetes' : str))
  .map(upstreamType => {
    return {
      key: upstreamType,
      value: _.startCase(_.toLower(upstreamType))
    };
  });

export enum UPSTREAM_SPEC_TYPES {
  AZURE = 'Azure',
  KUBE = 'Kubernetes',
  AWS = 'Aws',
  STATIC = 'Static',
  CONSUL = 'Consul'
}

export const CheckboxFilters: CheckboxFilterProps[] = Object.keys(
  UpstreamInput.SpecCase
)
  .filter(str => str !== 'SPEC_NOT_SET') // auto generated for oneof fields
  .map(str => (str === 'KUBE' ? 'Kubernetes' : str))
  .filter(str => str !== 'AWS_EC2')
  .map(str => {
    return {
      displayName: _.startCase(_.toLower(str)),
      value: false
    };
  });
// .concat({ displayName: 'other', value: false });

export const RadioFilters: RadioFilterProps = {
  options: [
    {
      displayName: 'Accepted'
    },
    {
      displayName: 'Rejected'
    }
  ]
};
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
