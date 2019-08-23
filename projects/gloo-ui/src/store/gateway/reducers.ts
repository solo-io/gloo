import { Gateway } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v2/gateway_pb';
import { GatewayActionTypes, GatewayAction } from './types';
import { GatewayDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';

export interface GatewayState {
  gatewaysList: GatewayDetails.AsObject[];
}

const initialState: GatewayState = {
  gatewaysList: []
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
    default:
      return state;
  }
}
