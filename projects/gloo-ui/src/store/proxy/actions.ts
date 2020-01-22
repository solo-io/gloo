import { Dispatch } from 'redux';
import { proxyAPI } from './api';
import { ListProxiesAction, ProxyAction } from './types';

export const listProxies = () => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await proxyAPI.getListProxies();
      dispatch<ListProxiesAction>({
        type: ProxyAction.LIST_PROXIES,
        payload: response
      });
    } catch (error) {}
  };
};
