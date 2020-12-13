import {
  DeleteUpstreamRequest,
  UpstreamDetails
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';

export enum UpstreamAction {
  GET_UPSTREAM = 'GET_UPSTREAM',
  LIST_UPSTREAMS = 'LIST_UPSTREAMS',
  CREATE_UPSTREAM = 'CREATE_UPSTREAM',
  CREATE_AWS_UPSTREAM = 'CREATE_AWS_UPSTREAM',
  UPDATE_UPSTREAM = 'UPDATE_UPSTREAM',
  DELETE_UPSTREAM = 'DELETE_UPSTREAM'
}

export interface ListUpstreamsAction {
  type: typeof UpstreamAction.LIST_UPSTREAMS;
  payload: UpstreamDetails.AsObject[];
}

export interface GetUpstreamAction {
  type: typeof UpstreamAction.GET_UPSTREAM;
  payload: UpstreamDetails.AsObject;
}
export interface CreateUpstreamAction {
  type: typeof UpstreamAction.CREATE_UPSTREAM;
  payload: UpstreamDetails.AsObject;
}

export interface UpdateUpstreamAction {
  type: typeof UpstreamAction.UPDATE_UPSTREAM;
  payload: UpstreamDetails.AsObject;
}

export interface DeleteUpstreamAction {
  type: typeof UpstreamAction.DELETE_UPSTREAM;
  payload: DeleteUpstreamRequest.AsObject;
}

export type UpstreamActionTypes =
  | ListUpstreamsAction
  | GetUpstreamAction
  | CreateUpstreamAction
  | UpdateUpstreamAction
  | DeleteUpstreamAction;
