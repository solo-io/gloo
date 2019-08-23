import { ProxyDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/proxy_pb';

export enum ProxyAction {
  GET_PROXY = 'GET_PROXY',
  LIST_PROXIES = 'LIST_PROXIES'
}

export interface ListProxiesAction {
  type: typeof ProxyAction.LIST_PROXIES;
  payload: ProxyDetails.AsObject[];
}

export interface GetProxyAction {
  type: typeof ProxyAction.GET_PROXY;
  payload: ProxyDetails.AsObject;
}

export type ProxyActionTypes = ListProxiesAction | GetProxyAction;
