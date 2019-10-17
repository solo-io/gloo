import {
  RouteTableDetails,
  DeleteRouteTableResponse,
  DeleteRouteTableRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/routetable_pb';

export enum RouteTableAction {
  GET_ROUTE_TABLE = 'GET_ROUTE_TABLE',
  LIST_ROUTE_TABLES = 'LIST_ROUTE_TABLES',
  CREATE_ROUTE_TABLE = 'CREATE_ROUTE_TABLE',
  UPDATE_ROUTE_TABLE_YAML = 'UPDATE_ROUTE_TABLE_YAML',
  UPDATE_ROUTE_TABLE_YAML_ERROR = 'UPDATE_ROUTE_TABLE_YAML_ERROR',
  UPDATE_ROUTE_TABLE = 'UPDATE_ROUTE_TABLE',
  DELETE_ROUTE_TABLE = 'DELETE_ROUTE_TABLE'
}

export interface ListRouteTablesAction {
  type: typeof RouteTableAction.LIST_ROUTE_TABLES;
  payload: RouteTableDetails.AsObject[];
}

export interface GetRouteTableAction {
  type: typeof RouteTableAction.GET_ROUTE_TABLE;
  payload: RouteTableDetails.AsObject;
}

export interface CreateRouteTableAction {
  type: typeof RouteTableAction.CREATE_ROUTE_TABLE;
  payload: RouteTableDetails.AsObject;
}

export interface UpdateRouteTableAction {
  type: typeof RouteTableAction.UPDATE_ROUTE_TABLE;
  payload: RouteTableDetails.AsObject;
}

export interface UpdateRouteTableYamlAction {
  type: typeof RouteTableAction.UPDATE_ROUTE_TABLE_YAML;
  payload: RouteTableDetails.AsObject;
}

export interface UpdateRouteTableYamlErrorAction {
  type: typeof RouteTableAction.UPDATE_ROUTE_TABLE_YAML_ERROR;
  payload: Error;
}

export interface DeleteRouteTableAction {
  type: typeof RouteTableAction.DELETE_ROUTE_TABLE;
  payload: DeleteRouteTableRequest.AsObject;
}

export type RouteTableActionTypes =
  | ListRouteTablesAction
  | GetRouteTableAction
  | CreateRouteTableAction
  | UpdateRouteTableAction
  | UpdateRouteTableYamlAction
  | UpdateRouteTableYamlErrorAction
  | DeleteRouteTableAction;
