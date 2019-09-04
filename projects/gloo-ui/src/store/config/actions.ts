import {
  GetVersionRequest,
  UpdateSettingsRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config_pb';
import { Dispatch } from 'redux';
import { showLoading, hideLoading } from 'react-redux-loading-bar';
import { config } from 'Api/v2/ConfigClient';
import {
  ConfigAction,
  GetVersionAction,
  GetSettingsAction,
  ListNamespacesAction,
  GetOAuthEndpointAction,
  GetIsLicenseValidAction,
  GetPodNamespaceAction,
  UpdateSettingsAction,
  UpdateWatchNamespacesAction,
  UpdateRefreshRateAction
} from './types';
import { Modal } from 'antd';
import { SuccessMessageAction, MessageAction } from 'store/modal/types';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
const { warning } = Modal;

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
        payload: response.toObject().settings!
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

export const updateWatchNamespaces = (updateWatchNamespacesRequest: {
  watchNamespacesList: string[];
}) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await config.updateWatchNamespaces(
        updateWatchNamespacesRequest
      );
      dispatch<UpdateWatchNamespacesAction>({
        type: ConfigAction.UPDATE_WATCH_NAMESPACES,
        payload: response.settings!
      });
      dispatch(hideLoading());
      dispatch<SuccessMessageAction>({
        type: MessageAction.SUCCESS_MESSAGE,
        message: 'Watched namespaces successfully updated.'
      });
    } catch (error) {
      warning({
        title: 'There was an error updating watched namespaces.',
        content: error.message
      });
    }
  };
};

export const updateRefreshRate = (updateRefreshRateRequest: {
  refreshRate: Duration.AsObject;
}) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await config.updateRefreshRate(updateRefreshRateRequest);
      dispatch<UpdateRefreshRateAction>({
        type: ConfigAction.UPDATE_REFRESH_RATE,
        payload: response.settings!
      });
      dispatch(hideLoading());
      dispatch<SuccessMessageAction>({
        type: MessageAction.SUCCESS_MESSAGE,
        message: 'Refresh rate successfully updated.'
      });
    } catch (error) {
      warning({
        title: 'There was an error updating refresh rate.',
        content: error.message
      });
    }
  };
};

export const updateSettings = (
  updateSettingsRequest: UpdateSettingsRequest
) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await config.updateSettings(updateSettingsRequest);
      dispatch<UpdateSettingsAction>({
        type: ConfigAction.UPDATE_SETTINGS,
        payload: response.settings!
      });
      dispatch(hideLoading());
      dispatch<SuccessMessageAction>({
        type: MessageAction.SUCCESS_MESSAGE,
        message: 'Settings successfully updated.'
      });
    } catch (error) {
      warning({
        title: 'There was an error updating settings.',
        content: error.message
      });
    }
  };
};
