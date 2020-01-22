import { Dispatch } from 'redux';
import { showLoading, hideLoading } from 'react-redux-loading-bar';
import { envoyAPI } from './api';
import { ListEnvoyDetailsAction, EnvoyAction } from './types';

export const listEnvoyDetails = () => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());
    try {
      const response = await envoyAPI.getEnvoyList();
      dispatch<ListEnvoyDetailsAction>({
        type: EnvoyAction.LIST_ENVOY_DETAILS,
        payload: response!
      });
      // dispatch(hideLoading());
    } catch (error) {}
  };
};
