import { OAuthEndpoint } from 'proto/solo-projects/projects/grpcserver/api/v1/config_pb';
import { Settings } from 'proto/gloo/projects/gloo/api/v1/settings_pb';
import { ConfigActionTypes, ConfigAction } from './types';

export interface ConfigState {
  version: string;
  oAuthEndpoint: OAuthEndpoint.AsObject;
  isLicenseValid: boolean;
  settings: Settings.AsObject;
  namespacesList: string[];
  namespace: string;
}

// TODO normalize
const initialState: ConfigState = {
  version: '',
  isLicenseValid: false,
  namespace: '',
  namespacesList: [] as string[],
  oAuthEndpoint: {
    clientName: '',
    url: ''
  },
  settings: ({} as unknown) as Settings.AsObject
};

export function configReducer(state = initialState, action: ConfigActionTypes) {
  switch (action.type) {
    case ConfigAction.GET_VERSION:
      return {
        ...state,
        version: action.payload
      };

    case ConfigAction.LIST_NAMESPACES:
      return {
        ...state,
        namespacesList: action.payload
      };

    case ConfigAction.GET_OAUTH_ENDPOINT:
      return {
        ...state,
        oAuthEndpoint: action.payload
      };

    case ConfigAction.GET_IS_LICENSE_VALID:
      return {
        ...state,
        isLicenseValid: action.payload
      };

    case ConfigAction.GET_POD_NAMESPACE:
      return {
        ...state,
        namespace: action.payload
      };

    case ConfigAction.GET_SETTINGS:
      return {
        ...state,
        settings: action.payload
      };

    default:
      return state;
  }
}
