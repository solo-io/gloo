import { Secret } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import { SecretActionTypes, SecretAction } from './types';

export interface SecretState {
  secretsList: Secret.AsObject[];
}

const initialState: SecretState = {
  secretsList: []
};

export function secretsReducer(
  state = initialState,
  action: SecretActionTypes
): SecretState {
  switch (action.type) {
    case SecretAction.LIST_SECRETS:
      return {
        ...state,
        secretsList: [...action.payload]
      };
    case SecretAction.CREATE_SECRET:
      return {
        ...state,
        secretsList: [...state.secretsList, action.payload]
      };
    case SecretAction.DELETE_SECRET:
      return {
        ...state,
        secretsList: state.secretsList.filter(
          s => s.metadata!.name !== action.payload.ref!.name
        )
      };
    default:
      return state;
  }
}
