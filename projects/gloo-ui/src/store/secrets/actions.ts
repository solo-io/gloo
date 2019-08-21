import {
  ListSecretsRequest,
  CreateSecretRequest,
  DeleteSecretRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { Dispatch } from 'redux';
import { showLoading, hideLoading } from 'react-redux-loading-bar';
import { secrets, getCreateSecret } from 'Api/v2/SecretClient';
import { ListSecretsAction, SecretAction, CreateSecretAction } from './types';

export const listSecrets = (
  listSecretsRequest: ListSecretsRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await secrets.getSecretsList(listSecretsRequest);
      dispatch<ListSecretsAction>({
        type: SecretAction.LIST_SECRETS,
        payload: response.secretsList
      });
      dispatch(hideLoading());
    } catch (error) {
      // handle error
    }
  };
};

export const createSecret = (
  createSecretRequest: CreateSecretRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await getCreateSecret(createSecretRequest);
      dispatch<CreateSecretAction>({
        type: SecretAction.CREATE_SECRET,
        payload: response.secret!
      });
    } catch (error) {
      //handle error
    }
  };
};

export const deleteSecret = (
  deleteSecretRequest: DeleteSecretRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await secrets.deleteSecret(deleteSecretRequest);
      dispatch<CreateSecretAction>({
        type: SecretAction.CREATE_SECRET,
        payload: response
      });
    } catch (error) {
      //handle error
    }
  };
};
