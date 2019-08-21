import { VirtualService } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { VirtualServiceActionTypes, VirtualServiceAction } from './types';
import { VirtualServiceDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { Route } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';

export interface VirtualServiceState {
  virtualServicesList: VirtualServiceDetails.AsObject[];
}

const initialState: VirtualServiceState = {
  virtualServicesList: []
};

export function virtualServicesReducer(
  state = initialState,
  action: VirtualServiceActionTypes
) {
  switch (action.type) {
    case VirtualServiceAction.LIST_VIRTUAL_SERVICES:
      return {
        ...state,
        virtualServicesList: [...action.payload]
      };
    case VirtualServiceAction.DELETE_VIRTUAL_SERVICE:
      return {
        ...state,
        virtualServicesList: state.virtualServicesList.filter(
          vs => vs.virtualService!.metadata!.name !== action.payload.ref!.name
        )
      };
    case VirtualServiceAction.DELETE_ROUTE:
      return {
        ...state
      };
    default:
      return state;
  }
}

export function routesReducer(
  state: VirtualServiceDetails.AsObject,
  action: VirtualServiceActionTypes
) {
  switch (action.type) {
    case VirtualServiceAction.DELETE_ROUTE:
      return { ...state };

    default:
      break;
  }
}
