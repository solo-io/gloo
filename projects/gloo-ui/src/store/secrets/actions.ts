import {
  ListSecretsRequest,
  CreateSecretRequest,
  DeleteSecretRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { Dispatch } from 'redux';
import { showLoading, hideLoading } from 'react-redux-loading-bar';
import { secrets } from './api';
import { ListSecretsAction, SecretAction, CreateSecretAction } from './types';
import { Modal } from 'antd';
import { guardByLicense } from 'store/config/actions';
import { SoloWarning } from 'Components/Common/SoloWarningContent';
const { warning } = Modal;

export const listSecrets = () => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());
    try {
      const response = await secrets.getSecretsList();
      dispatch<ListSecretsAction>({
        type: SecretAction.LIST_SECRETS,
        payload: response.secretsList
      });
      // dispatch(hideLoading());
    } catch (error) {
      // handle error
    }
  };
};

export const createSecret = (
  createSecretRequest: CreateSecretRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());
    try {
      const response = await secrets.createSecret(createSecretRequest);
      dispatch<CreateSecretAction>({
        type: SecretAction.CREATE_SECRET,
        payload: response.secret!
      });
    } catch (error) {
      SoloWarning('There was an error creating the secret.', error);
    }
  };
};

export const deleteSecret = (
  deleteSecretRequest: DeleteSecretRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());
    try {
      guardByLicense();
      const response = await secrets.deleteSecret(deleteSecretRequest);
      dispatch<CreateSecretAction>({
        type: SecretAction.CREATE_SECRET,
        payload: response
      });
    } catch (error) {
      SoloWarning('There was an error deleting the secret.', error);
    }
  };
};
