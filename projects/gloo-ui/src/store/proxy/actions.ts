import { Dispatch } from 'redux';
import { proxys } from './api';
import { ListProxiesAction, ProxyAction } from './types';

export const listProxies = () => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await proxys.getListProxies();
      dispatch<ListProxiesAction>({
        type: ProxyAction.LIST_PROXIES,
        payload: response.proxyDetailsList
      });
    } catch (error) {}
  };
};
