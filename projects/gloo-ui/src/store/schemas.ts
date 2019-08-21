import { normalize, schema } from 'normalizr';

export const upstream = new schema.Entity(
  'upstreams',
  {},
  {
    idAttribute: (value, parent, key) => {
      return `${value.metadata!.name}-${value.metadata!.namespace}`;
    }
  }
);

export const route = new schema.Entity('routes');

export const virtualService = new schema.Entity('virtualServices', {
  routesList: [route]
});

/*
// export namespace VirtualServiceDetails {
//     export type AsObject = {
//          virtualService?: {
//             virtualHost ?: {
//               name: string,
//               domainsList: Array < string >,
//               routesList: Array < {
                        matcher?: {
                            prefix: string,
                            exact: string,
                            regex: string,
                            headersList: Array<HeaderMatcher.AsObject>,
                            queryParametersList: Array<QueryParameterMatcher.AsObject>,
                            methodsList: Array<string>,
                        }
                        routeAction?: {
                            single?:  {
                                upstream?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
                                kube?: KubernetesServiceDestination.AsObject,
                                consul?: ConsulServiceDestination.AsObject,
                                destinationSpec?:  {
                                    aws?: {
                                        logicalName: string,
                                        invocationStyle: DestinationSpec.InvocationStyle,
                                        responseTransformation: boolean,
                                    },
                                    azure?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_azure_azure_pb.DestinationSpec.AsObject,
                                    rest?: {
                                        functionName: string,
                                        parameters?:  {
                                            headersMap: Array<[string, string]>,
                                            path?: google_protobuf_wrappers_pb.StringValue.AsObject,
                                        },
                                        responseTransformation?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_transformation_transformation_pb.TransformationTemplate.AsObject,
                                    }
                                    grpc?: {
                                        pb_package: string,
                                        service: string,
                                        pb_function: string,
                                        parameters?: {
                                            headersMap: Array<[string, string]>,
                                            path?: google_protobuf_wrappers_pb.StringValue.AsObject,
                                        },
                                    },
                                },
                                subset?: github_com_solo_io_gloo_projects_gloo_api_v1_subset_pb.Subset.AsObject,
                            }
                            multi?: MultiDestination.AsObject,
                            upstreamGroup?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
                        }
                        redirectAction?: RedirectAction.AsObject,
                        directResponseAction?: DirectResponseAction.AsObject,
                        routePlugins?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.RoutePlugins.AsObject,
  }>,
//               virtualHostPlugins ?: github_com_solo_io_gloo_projects_gloo_api_v1_plugins_pb.VirtualHostPlugins.AsObject,
//               corsPolicy ?: CorsPolicy.AsObject,
//   }
    //         sslConfig ?: github_com_solo_io_gloo_projects_gloo_api_v1_ssl_pb.SslConfig.AsObject,
    //         displayName: string,
//             status ?: {
//               state: Status.State,
//               reason: string,
//               reportedBy: string,
//               subresourceStatusesMap: Array < [string, Status.AsObject] >,
//             metadata ?: {
//               name: string,
//               namespace: string,
//               cluster: string,
//               resourceVersion: string,
//               labelsMap: Array < [string, string] >,
//               annotationsMap: Array < [string, string] >,
//   },
//             }
//         plugins?: Plugins.AsObject,
//         raw?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw.AsObject,
//     }
// }
*/
