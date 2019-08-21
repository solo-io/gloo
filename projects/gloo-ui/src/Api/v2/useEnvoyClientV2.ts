import { envoy } from './EnvoyClient';
import { useSendRequest } from './requestReducerV2';
import { ListEnvoyDetailsRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/envoy_pb';

export function useGetEnvoyList(variables?: ListEnvoyDetailsRequest.AsObject) {
  return useSendRequest(variables, envoy.getEnvoyList);
}
