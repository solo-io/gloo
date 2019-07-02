import * as React from 'react';
import { requestReducer, Reducer, RequestAction } from './request-reducer';

import { secretClient } from './grpc-web-hooks';
import {
  GetSecretRequest,
  GetSecretResponse,
  ListSecretsRequest,
  ListSecretsResponse,
  CreateSecretRequest,
  CreateSecretResponse,
  UpdateSecretRequest,
  UpdateSecretResponse,
  DeleteSecretRequest,
  DeleteSecretResponse
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb_service';

/* -------------------------------------------------------------------------- */
/*                                 GET SECRET                                 */
/* -------------------------------------------------------------------------- */

export const useGetSecret = (
  request: GetSecretRequest | null,
  initialData: GetSecretResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<GetSecretResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: GetSecretRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    secretClient.getSecret(
      request,
      (
        error: ServiceError | null,
        responseMessage: GetSecretResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};

/* -------------------------------------------------------------------------- */
/*                                LIST SECRETS                                */
/* -------------------------------------------------------------------------- */

export const useListSecrets = (
  request: ListSecretsRequest | null,
  initialData: ListSecretsResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<ListSecretsResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: ListSecretsRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    secretClient.listSecrets(
      request,
      (
        error: ServiceError | null,
        responseMessage: ListSecretsResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
/* -------------------------------------------------------------------------- */
/*                                CREATE SECRET                               */
/* -------------------------------------------------------------------------- */

export const useCreateSecret = (
  request: CreateSecretRequest | null,
  initialData: CreateSecretResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<CreateSecretResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: CreateSecretRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    secretClient.createSecret(
      request,
      (
        error: ServiceError | null,
        responseMessage: CreateSecretResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
/* -------------------------------------------------------------------------- */
/*                                UPDATE SECRET                               */
/* -------------------------------------------------------------------------- */

export const useUpdateSecret = (
  request: UpdateSecretRequest | null,
  initialData: UpdateSecretResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<UpdateSecretResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: UpdateSecretRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    secretClient.updateSecret(
      request,
      (
        error: ServiceError | null,
        responseMessage: UpdateSecretResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
/* -------------------------------------------------------------------------- */
/*                                DELETE SECRET                               */
/* -------------------------------------------------------------------------- */

export const useDeleteSecret = (
  request: DeleteSecretRequest | null,
  initialData: DeleteSecretResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<DeleteSecretResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: DeleteSecretRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    secretClient.deleteSecret(
      request,
      (
        error: ServiceError | null,
        responseMessage: DeleteSecretResponse | null
      ) => {
        if (error) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          if (!mounted.current) return;
          dispatch({ type: RequestAction.ERROR, payload: null, error });
        } else {
          const response = responseMessage;
          if (!mounted.current) return;
          dispatch({
            type: RequestAction.SUCCESS,
            payload: response!.toObject()
          });
        }
      }
    );
  };

  React.useEffect(() => {
    makeRequest(request || null);
    return () => {
      mounted.current = false;
    };
  }, []);

  return {
    data: state.data!,
    error: state.error,
    loading: state.isLoading,
    refetch: makeRequest
  };
};
