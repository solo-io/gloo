import {
  GetVersionRequest,
  UpdateSettingsRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config_pb';
import { Dispatch } from 'redux';
import { showLoading } from 'react-redux-loading-bar';
import { config } from 'Api/v2/ConfigClient';
import {
  ConfigAction,
  GetVersionAction,
  GetSettingsAction,
  ListNamespacesAction,
  GetOAuthEndpointAction,
  GetIsLicenseValidAction,
  GetPodNamespaceAction,
  UpdateSettingsAction
} from './types';

export const getVersion = () => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await config.getVersion();
      dispatch<GetVersionAction>({
        type: ConfigAction.GET_VERSION,
        payload: response.version
      });
    } catch (error) {}
  };
};
export const getSettings = () => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await config.getSettings();
      dispatch<GetSettingsAction>({
        type: ConfigAction.GET_SETTINGS,
        payload: response.settings!
      });
    } catch (error) {}
  };
};
export const listNamespaces = () => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await config.listNamespaces();
      dispatch<ListNamespacesAction>({
        type: ConfigAction.LIST_NAMESPACES,
        payload: response.namespacesList!
      });
    } catch (error) {}
  };
};

export const getOAuthEndpoint = () => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await config.getOAuthEndpoint();
      dispatch<GetOAuthEndpointAction>({
        type: ConfigAction.GET_OAUTH_ENDPOINT,
        payload: response.oAuthEndpoint!
      });
    } catch (error) {}
  };
};

export const getIsLicenseValid = () => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await config.getIsLicenseValid();
      dispatch<GetIsLicenseValidAction>({
        type: ConfigAction.GET_IS_LICENSE_VALID,
        payload: response.isLicenseValid!
      });
    } catch (error) {}
  };
};
export const getPodNamespace = () => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await config.getPodNamespace();
      dispatch<GetPodNamespaceAction>({
        type: ConfigAction.GET_POD_NAMESPACE,
        payload: response.namespace
      });
    } catch (error) {}
  };
};
export const updateSettings = (
  updateSettingsRequest: UpdateSettingsRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await config.updateSettings(updateSettingsRequest);
      dispatch<UpdateSettingsAction>({
        type: ConfigAction.UPDATE_SETTINGS,
        payload: response.settings!
      });
    } catch (error) {}
  };
};
