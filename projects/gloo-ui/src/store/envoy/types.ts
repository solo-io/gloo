import { EnvoyDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/envoy_pb';

export enum EnvoyAction {
  LIST_ENVOY_DETAILS = 'LIST_ENVOY_DETAILS'
}

export interface ListEnvoyDetailsAction {
  type: typeof EnvoyAction.LIST_ENVOY_DETAILS;
  payload: EnvoyDetails.AsObject[];
}

export type EnvoyActionTypes = ListEnvoyDetailsAction;
