import { gateways } from './GatewayClient';
import { useSendRequest } from './requestReducerV2';

export function useGetGatewayList(variables: {namespaces: string[] }) {
  return useSendRequest(variables, gateways.getGatewaysList);
}
