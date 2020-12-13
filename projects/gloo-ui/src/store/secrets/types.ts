import { Secret } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import { DeleteSecretRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';

export enum SecretAction {
  GET_SECRET = 'GET_SECRET',
  LIST_SECRETS = 'LIST_SECRETS',
  UPDATE_SECRET = 'UPDATE_SECRET',
  CREATE_SECRET = 'CREATE_SECRET',
  DELETE_SECRET = 'DELETE_SECRET'
}

export interface ListSecretsAction {
  type: typeof SecretAction.LIST_SECRETS;
  payload: Secret.AsObject[];
}

export interface GetSecretAction {
  type: typeof SecretAction.GET_SECRET;
  payload: Secret.AsObject;
}

export interface CreateSecretAction {
  type: typeof SecretAction.CREATE_SECRET;
  payload: Secret.AsObject;
}

export interface DeleteSecretAction {
  type: typeof SecretAction.DELETE_SECRET;
  payload: DeleteSecretRequest.AsObject;
}

export interface UpdateSecretAction {
  type: typeof SecretAction.UPDATE_SECRET;
  payload: Secret.AsObject;
}

export type SecretActionTypes =
  | ListSecretsAction
  | GetSecretAction
  | CreateSecretAction
  | DeleteSecretAction
  | UpdateSecretAction;
