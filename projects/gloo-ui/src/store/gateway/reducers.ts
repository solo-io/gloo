import { GatewayActionTypes, GatewayAction } from './types';
import { GatewayDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import { SoloWarning } from 'Components/Common/SoloWarningContent';

export interface GatewayState {
  gatewaysList: GatewayDetails.AsObject[];
  yamlParseError: boolean;
}

const initialState: GatewayState = {
  gatewaysList: [],
  yamlParseError: false
};

export function gatewaysReducer(
  state = initialState,
  action: GatewayActionTypes
): GatewayState {
  switch (action.type) {
    case GatewayAction.GET_GATEWAY:
      return {
        ...state,
        gatewaysList: [...state.gatewaysList, action.payload]
      };
    case GatewayAction.LIST_GATEWAYS:
      return { ...state, gatewaysList: [...action.payload] };
    case GatewayAction.UPDATE_GATEWAY:
      return {
        ...state
      };
    case GatewayAction.UPDATE_GATEWAY_YAML_ERROR:
      SoloWarning(
        'There was an error updating the virtual service.',
        action.payload
      );
      return {
        ...state,
        yamlParseError: true
      };
    case GatewayAction.UPDATE_GATEWAY_YAML:
      return {
        ...state,
        yamlParseError: false
      };
    default:
      return { ...state, yamlParseError: false };
  }
}
