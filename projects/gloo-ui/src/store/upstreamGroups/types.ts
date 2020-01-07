import {
  UpstreamGroupDetails,
  DeleteUpstreamGroupResponse
} from 'proto/solo-projects/projects/grpcserver/api/v1/upstreamgroup_pb';
import { UpdateVirtualServiceYamlErrorAction } from 'store/virtualServices/types';

export enum UpstreamGroupAction {
  GET_UPSTREAM_GROUP = 'GET_UPSTREAM_GROUP',
  LIST_UPSTREAM_GROUPS = 'LIST_UPSTREAM_GROUPS',
  CREATE_UPSTREAM_GROUP = 'CREATE_UPSTREAM_GROUP',
  UPDATE_UPSTREAM_GROUP = 'UPDATE_UPSTREAM_GROUP',
  UPDATE_UPSTREAM_GROUP_YAML = 'UPDATE_UPSTREAM_GROUPYAML',
  UPDATE_UPSTREAM_GROUP_YAML_ERROR = 'UPDATE_UPSTREAM_GROUP_YAML_ERROR',

  DELETE_UPSTREAM_GROUP = 'DELETE_UPSTREAM_GROUP'
}

export interface GetUpstreamGroupAction {
  type: typeof UpstreamGroupAction.GET_UPSTREAM_GROUP;
  payload: UpstreamGroupDetails;
}

export interface ListUpstreamGroupsAction {
  type: typeof UpstreamGroupAction.LIST_UPSTREAM_GROUPS;
  payload: UpstreamGroupDetails.AsObject[];
}

export interface CreateUpstreamGroupAction {
  type: typeof UpstreamGroupAction.CREATE_UPSTREAM_GROUP;
  payload: UpstreamGroupDetails.AsObject;
}

export interface UpdateUpstreamGroupAction {
  type: typeof UpstreamGroupAction.UPDATE_UPSTREAM_GROUP;
  payload: UpstreamGroupDetails.AsObject;
}

export interface DeleteUpstreamGroupAction {
  type: typeof UpstreamGroupAction.DELETE_UPSTREAM_GROUP;
  payload: DeleteUpstreamGroupResponse.AsObject;
}

export interface UpdateUpstreamGroupYamlAction {
  type: typeof UpstreamGroupAction.UPDATE_UPSTREAM_GROUP_YAML;
  payload: UpstreamGroupDetails.AsObject;
}
export interface UpdateUpstreamGroupYamlError {
  type: typeof UpstreamGroupAction.UPDATE_UPSTREAM_GROUP_YAML_ERROR;
  payload: Error;
}

export type UpstreamGroupActionTypes =
  | GetUpstreamGroupAction
  | ListUpstreamGroupsAction
  | CreateUpstreamGroupAction
  | UpdateUpstreamGroupAction
  | DeleteUpstreamGroupAction
  | UpdateUpstreamGroupYamlAction
  | UpdateUpstreamGroupYamlError;
