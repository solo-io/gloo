import { ProxyDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/proxy_pb';
import { ProxyActionTypes, ProxyAction } from './types';

export interface ProxyState {
  proxiesList: ProxyDetails.AsObject[];
}

const initialState: ProxyState = {
  proxiesList: []
};

export function proxyReducer(
  state = initialState,
  action: ProxyActionTypes
): ProxyState {
  switch (action.type) {
    case ProxyAction.LIST_PROXIES:
      return {
        ...state,
        proxiesList: [...action.payload]
      };

    case ProxyAction.GET_PROXY:
      return {
        ...state
      };
    default:
      return state;
  }
}
