import { EnvoyDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/envoy_pb';
import { EnvoyActionTypes, EnvoyAction } from './types';

export interface EnvoyState {
  envoyDetailsList: EnvoyDetails.AsObject[];
}

const initialState: EnvoyState = {
  envoyDetailsList: []
};
export function envoyReducer(state = initialState, action: EnvoyActionTypes) {
  switch (action.type) {
    case EnvoyAction.LIST_ENVOY_DETAILS:
      return {
        envoyDetailsList: [...action.payload]
      };

    default:
      return state;
  }
}
