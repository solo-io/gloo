import { VirtualServiceActionTypes, VirtualServiceAction } from './types';
import { VirtualServiceDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { SoloWarning } from 'Components/Common/SoloWarningContent';

export interface VirtualServiceState {
  virtualServicesList: VirtualServiceDetails.AsObject[];
  yamlParseError: boolean;
}

const initialState: VirtualServiceState = {
  virtualServicesList: [],
  yamlParseError: false
};

export function virtualServicesReducer(
  state = initialState,
  action: VirtualServiceActionTypes
): VirtualServiceState {
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
    case VirtualServiceAction.UPDATE_VIRTUAL_SERVICE_YAML_ERROR:
      SoloWarning(
        'There was an error updating the virtual service.',
        action.payload
      );
      return {
        ...state,
        yamlParseError: true
      };
    case VirtualServiceAction.UPDATE_VIRTUAL_SERVICE_YAML:
      return {
        ...state,
        yamlParseError: false
      };
    default:
      return { ...state, yamlParseError: false };
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
