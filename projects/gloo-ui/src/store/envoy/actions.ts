import { Dispatch } from 'redux';
import { showLoading, hideLoading } from 'react-redux-loading-bar';
import { envoy } from 'Api/v2/EnvoyClient';
import { ListEnvoyDetailsAction, EnvoyAction } from './types';
import { Modal } from 'antd';
const { warning } = Modal;

export const listEnvoyDetails = () => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());
    try {
      const response = await envoy.getEnvoyList();
      dispatch<ListEnvoyDetailsAction>({
        type: EnvoyAction.LIST_ENVOY_DETAILS,
        payload: response!.toObject().envoyDetailsList
      });
      // dispatch(hideLoading());
    } catch (error) {}
  };
};
