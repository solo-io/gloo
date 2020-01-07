import { UpstreamGroupDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/upstreamgroup_pb';
import { UpstreamGroupActionTypes, UpstreamGroupAction } from './types';

export interface UpstreamGroupState {
  upstreamGroupsList: UpstreamGroupDetails.AsObject[];
  yamlParseError: boolean;
}

const initialState: UpstreamGroupState = {
  upstreamGroupsList: [],
  yamlParseError: false
};

export function upstreamGroupsReducer(
  state = initialState,
  action: UpstreamGroupActionTypes
): UpstreamGroupState {
  switch (action.type) {
    case UpstreamGroupAction.LIST_UPSTREAM_GROUPS:
      return {
        ...state,
        upstreamGroupsList: action.payload
      };
    case UpstreamGroupAction.UPDATE_UPSTREAM_GROUP_YAML_ERROR:
      return {
        ...state,
        yamlParseError: true
      };

    default:
      break;
  }
  return { ...state, yamlParseError: false };
}
