import { secrets } from './SecretClient';
import { useSendRequest } from './requestReducerV2';

export function useGetSecretsListV2(variables: { namespaces: string[] }) {
  return useSendRequest(variables, secrets.getSecretsList);
}
