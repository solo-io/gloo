import { proxy } from './ProxyClient';
import { useSendRequest } from './requestReducerV2';

export function useGetProxiesList(variables: { namespaces: string[] }) {
  return useSendRequest(variables, proxy.getProxiesList);
}
