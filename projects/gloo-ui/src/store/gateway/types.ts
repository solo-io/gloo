import { GatewayDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';

export enum GatewayAction {
  GET_GATEWAY = 'GET_GATEWAY',
  LIST_GATEWAYS = 'LIST_GATEWAYS',
  UPDATE_GATEWAY = 'UPDATE_GATEWAY'
}

export interface GetGatewayAction {
  type: typeof GatewayAction.GET_GATEWAY;
  payload: GatewayDetails.AsObject;
}
export interface ListGatewaysAction {
  type: typeof GatewayAction.LIST_GATEWAYS;
  payload: GatewayDetails.AsObject[];
}

export interface UpdateGatewayAction {
  type: typeof GatewayAction.UPDATE_GATEWAY;
  payload: GatewayDetails.AsObject;
}
export type GatewayActionTypes =
  | GetGatewayAction
  | ListGatewaysAction
  | UpdateGatewayAction;
