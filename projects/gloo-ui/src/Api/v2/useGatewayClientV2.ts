import { gateways, UpdateGatewayHttpData } from './GatewayClient';
import { useSendRequest } from './requestReducerV2';
import { UpdateGatewayRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import { Gateway } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v2/gateway_pb';

export function useGetGatewayList(variables: { namespaces: string[] }) {
  return useSendRequest(variables, gateways.getGatewaysList);
}

export function useUpdateGateway(
  variables: {
    originalGateway: Gateway;
    updates: UpdateGatewayHttpData;
  } | null
) {
  return useSendRequest(variables, gateways.updateGateway);
}
