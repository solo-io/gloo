import {
  ListProxiesRequest,
  ListProxiesResponse
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/proxy_pb';
import { Dispatch } from 'redux';
import { client } from 'Api/v2/ProxyClient';
import { ListProxiesAction, ProxyAction } from './types';

export function getListProxies(
  listProxiesRequest: ListProxiesRequest.AsObject
): Promise<ListProxiesResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new ListProxiesRequest();
    request.setNamespacesList(listProxiesRequest.namespacesList);
    client.listProxies(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        // TODO: normalize
        resolve(data!.toObject());
      }
    });
  });
}

export const listProxies = (
  listProxiesRequest: ListProxiesRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await getListProxies(listProxiesRequest);
      dispatch<ListProxiesAction>({
        type: ProxyAction.LIST_PROXIES,
        payload: response.proxyDetailsList
      });
    } catch (error) {}
  };
};
