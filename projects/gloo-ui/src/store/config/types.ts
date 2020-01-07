import {
  GetSettingsRequest,
  OAuthEndpoint
} from 'proto/solo-projects/projects/grpcserver/api/v1/config_pb';
import { Settings } from 'proto/gloo/projects/gloo/api/v1/settings_pb';

export enum ConfigAction {
  GET_VERSION = 'GET_VERSION',
  GET_OAUTH_ENDPOINT = 'GET_OAUTH_ENDPOINT',
  GET_IS_LICENSE_VALID = 'GET_IS_LICENSE_VALID',
  GET_SETTINGS = 'GET_SETTINGS',
  UPDATE_SETTINGS = 'UPDATE_SETTINGS',
  UPDATE_WATCH_NAMESPACES = 'UPDATE_WATCH_NAMESPACES',
  UPDATE_REFRESH_RATE = 'UPDATE_REFRESH_RATE',
  LIST_NAMESPACES = 'LIST_NAMESPACES',
  GET_POD_NAMESPACE = 'GET_POD_NAMESPACE',
  LICENSE_CHECK_GUARD = 'LICENSE_CHECK_GUARD'
}

export interface GetSettingsAction {
  type: typeof ConfigAction.GET_SETTINGS;
  payload: Settings.AsObject;
}

export interface ListNamespacesAction {
  type: typeof ConfigAction.LIST_NAMESPACES;
  payload: string[];
}

export interface GetOAuthEndpointAction {
  type: typeof ConfigAction.GET_OAUTH_ENDPOINT;
  payload: OAuthEndpoint.AsObject;
}

export interface GetIsLicenseValidAction {
  type: typeof ConfigAction.GET_IS_LICENSE_VALID;
  payload: boolean;
}

export interface GetVersionAction {
  type: typeof ConfigAction.GET_VERSION;
  payload: string;
}

export interface UpdateSettingsAction {
  type: typeof ConfigAction.UPDATE_SETTINGS;
  payload: Settings.AsObject;
}

export interface UpdateWatchNamespacesAction {
  type: typeof ConfigAction.UPDATE_WATCH_NAMESPACES;
  payload: Settings.AsObject;
}

export interface UpdateRefreshRateAction {
  type: typeof ConfigAction.UPDATE_REFRESH_RATE;
  payload: Settings.AsObject;
}

export interface GetPodNamespaceAction {
  type: typeof ConfigAction.GET_POD_NAMESPACE;
  payload: string;
}

export interface LicenseCheckGuard {
  type: typeof ConfigAction.LICENSE_CHECK_GUARD;
}

export type ConfigActionTypes =
  | GetSettingsAction
  | ListNamespacesAction
  | GetOAuthEndpointAction
  | GetIsLicenseValidAction
  | GetVersionAction
  | UpdateSettingsAction
  | UpdateWatchNamespacesAction
  | UpdateRefreshRateAction
  | GetPodNamespaceAction
  | LicenseCheckGuard;
