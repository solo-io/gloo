import { UpstreamDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { UpstreamAction, UpstreamActionTypes } from './types';

export interface UpstreamState {
  upstreamsList: UpstreamDetails.AsObject[];
}

/* 
export namespace UpstreamInput {
  export type AsObject = {
    ref?: {
    name: string,
    namespace: string,
  },
    kube?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_kubernetes_kubernetes_pb.UpstreamSpec.AsObject,
    pb_static?: {
    hostsList: Array<{
    addr: string,
    port: number,
  }>,
    useTls: boolean,
    serviceSpec?: {
    rest?: {
    transformationsMap: Array<[string, github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_transformation_pb.TransformationTemplate.AsObject]>,
    swaggerInfo?: {
      url: string,
      inline: string,
    },
  },
    grpc?: {
    descriptors: Uint8Array | string,
    grpcServicesList: Array<ServiceSpec.GrpcService.AsObject>,
  },
  },
  },
    aws?: {
    region: string,
    secretRef?: {
    name: string,
    namespace: string,
  },
    lambdaFunctionsList: {
    logicalName: string,
    lambdaFunctionName: string,
    qualifier: string,
  }[],
  },
    azure?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb.UpstreamSpec.AsObject,
    consul?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_consul_consul_pb.UpstreamSpec.AsObject,
    awsEc2?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_aws_ec2_aws_ec2_pb.UpstreamSpec.AsObject,
  }
*/
// add isLoading
// add error
/**
 * initial state = {
 * upstreams: {
 * aws: [],
 * azure: [],
 * static:[],
 * consul: [],
 * }}
 */
const initialState: UpstreamState = {
  upstreamsList: []
};

export function upstreamsReducer(
  state = initialState,
  action: UpstreamActionTypes
): UpstreamState {
  switch (action.type) {
    case UpstreamAction.LIST_UPSTREAMS:
      return { ...state, upstreamsList: [...action.payload] };

    case UpstreamAction.UPDATE_UPSTREAM:
      return {
        ...state,
        upstreamsList: state.upstreamsList.map(upstream =>
          upstream.upstream!.metadata!.name !==
          action.payload.upstream!.metadata!.name
            ? upstream
            : action.payload
        )
      };
    case UpstreamAction.CREATE_UPSTREAM:
      return {
        ...state,
        upstreamsList: [...state.upstreamsList, action.payload]
      };
    case UpstreamAction.DELETE_UPSTREAM:
      return {
        ...state,
        upstreamsList: state.upstreamsList.filter(
          upstream =>
            upstream.upstream!.metadata!.name !== action.payload.ref!.name
        )
      };
    default:
      return state;
  }
}
