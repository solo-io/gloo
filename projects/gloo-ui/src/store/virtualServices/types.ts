import {
  VirtualServiceDetails,
  DeleteVirtualServiceRequest
} from 'proto/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';

export enum VirtualServiceAction {
  GET_VIRTUAL_SERVICE = 'GET_VIRTUAL_SERVICE',
  LIST_VIRTUAL_SERVICES = 'LIST_VIRTUAL_SERVICES',
  CREATE_VIRTUAL_SERVICE = 'CREATE_VIRTUAL_SERVICE',
  UPDATE_VIRTUAL_SERVICE = 'UPDATE_VIRTUAL_SERVICE',
  DELETE_VIRTUAL_SERVICE = 'DELETE_VIRTUAL_SERVICE',
  UPDATE_VIRTUAL_SERVICE_YAML = 'UPDATE_VIRTUAL_SERVICE_YAML',
  UPDATE_VIRTUAL_SERVICE_YAML_ERROR = 'UPDATE_VIRTUAL_SERVICE_YAML_ERROR',
  CREATE_ROUTE = 'CREATE_ROUTE',
  UPDATE_ROUTE = 'UPDATE_ROUTE',
  DELETE_ROUTE = 'DELETE_ROUTE',
  SWAP_ROUTES = 'SWAP_ROUTES',
  SHIFT_ROUTES = 'SHIFT_ROUTES'
}

export interface GetVirtualServiceAction {
  type: typeof VirtualServiceAction.GET_VIRTUAL_SERVICE;
  payload: VirtualServiceDetails.AsObject;
}

export interface ListVirtualServicesAction {
  type: typeof VirtualServiceAction.LIST_VIRTUAL_SERVICES;
  payload: VirtualServiceDetails.AsObject[];
}

export interface CreateVirtualServiceAction {
  type: typeof VirtualServiceAction.CREATE_VIRTUAL_SERVICE;
  payload: VirtualServiceDetails.AsObject;
}
export interface UpdateVirtualServiceAction {
  type: typeof VirtualServiceAction.UPDATE_VIRTUAL_SERVICE;
  payload: VirtualServiceDetails.AsObject;
}
export interface DeleteVirtualServiceAction {
  type: typeof VirtualServiceAction.DELETE_VIRTUAL_SERVICE;
  payload: DeleteVirtualServiceRequest.AsObject;
}

export interface UpdateVirtualServiceYamlAction {
  type: typeof VirtualServiceAction.UPDATE_VIRTUAL_SERVICE_YAML;
  payload: VirtualServiceDetails.AsObject;
}

export interface UpdateVirtualServiceYamlErrorAction {
  type: typeof VirtualServiceAction.UPDATE_VIRTUAL_SERVICE_YAML_ERROR;
  payload: Error;
}

export interface CreateRouteAction {
  type: typeof VirtualServiceAction.CREATE_ROUTE;
  payload: VirtualServiceDetails.AsObject;
}

export interface UpdateRouteAction {
  type: typeof VirtualServiceAction.UPDATE_ROUTE;
  payload: VirtualServiceDetails.AsObject;
}
export interface DeleteRouteAction {
  type: typeof VirtualServiceAction.DELETE_ROUTE;
  payload: VirtualServiceDetails.AsObject;
}

// Routes
export interface SwapRouteAction {
  type: typeof VirtualServiceAction.SWAP_ROUTES;
  payload: VirtualServiceDetails.AsObject;
}
export interface ShiftRoutesAction {
  type: typeof VirtualServiceAction.SHIFT_ROUTES;
  payload: VirtualServiceDetails.AsObject;
}

export type VirtualServiceActionTypes =
  | GetVirtualServiceAction
  | ListVirtualServicesAction
  | CreateVirtualServiceAction
  | UpdateVirtualServiceAction
  | DeleteVirtualServiceAction
  | UpdateVirtualServiceYamlAction
  | UpdateVirtualServiceYamlErrorAction
  | CreateRouteAction
  | UpdateRouteAction
  | DeleteRouteAction
  | SwapRouteAction
  | ShiftRoutesAction;
