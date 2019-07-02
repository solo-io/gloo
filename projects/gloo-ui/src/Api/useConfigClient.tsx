import * as React from 'react';
import { requestReducer, Reducer, RequestAction } from './request-reducer';

import {
  GetVersionRequest,
  GetVersionResponse,
  GetOAuthEndpointRequest,
  GetOAuthEndpointResponse,
  GetIsLicenseValidRequest,
  GetIsLicenseValidResponse,
  GetSettingsRequest,
  GetSettingsResponse,
  UpdateSettingsRequest,
  UpdateSettingsResponse,
  ListNamespacesRequest,
  ListNamespacesResponse
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config_pb';
import { configClient } from './grpc-web-hooks';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config_pb_service';

/* -------------------------------------------------------------------------- */
/*                                 GET VERSION                                */
/* -------------------------------------------------------------------------- */

export const useGetVersion = (
  request: GetVersionRequest | null,
  initialData: GetVersionResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<GetVersionResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = () => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    configClient.getVersion(
      request,
      (
        error: ServiceError | null,
        responseMessage: GetVersionResponse | null
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
    makeRequest();
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
/*                             GET OAUTH ENDPOINT                             */
/* -------------------------------------------------------------------------- */

export const useGetOAuthEndpoint = (
  request: GetOAuthEndpointRequest | null,
  initialData: GetOAuthEndpointResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<GetOAuthEndpointResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = () => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    configClient.getOAuthEndpoint(
      request,
      (
        error: ServiceError | null,
        responseMessage: GetOAuthEndpointResponse | null
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
    makeRequest();
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
/*                            GET IS LICENSE VALID                            */
/* -------------------------------------------------------------------------- */

export const useGetIsLicenseValid = (
  request: GetIsLicenseValidRequest | null,
  initialData: GetIsLicenseValidResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<GetIsLicenseValidResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = () => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    configClient.getIsLicenseValid(
      request,
      (
        error: ServiceError | null,
        responseMessage: GetIsLicenseValidResponse | null
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
    makeRequest();
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
/*                                GET SETTINGS                                */
/* -------------------------------------------------------------------------- */

export const useGetSettings = (
  request: GetSettingsRequest | null,
  initialData: GetSettingsResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<GetSettingsResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = () => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    configClient.getSettings(
      request,
      (
        error: ServiceError | null,
        responseMessage: GetSettingsResponse | null
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
    makeRequest();
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
/*                               UPDATE SETTINGS                              */
/* -------------------------------------------------------------------------- */

export const useUpdateSettings = (
  request: UpdateSettingsRequest | null,
  initialData: UpdateSettingsResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<UpdateSettingsResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = () => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    configClient.updateSettings(
      request,
      (
        error: ServiceError | null,
        responseMessage: UpdateSettingsResponse | null
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
    makeRequest();
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
/*                               LIST NAMESPACES                              */
/* -------------------------------------------------------------------------- */

export const useListNamespaces = (
  request: ListNamespacesRequest | null,
  initialData: ListNamespacesResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<ListNamespacesResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = () => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    configClient.listNamespaces(
      request,
      (
        error: ServiceError | null,
        responseMessage: ListNamespacesResponse | null
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
    makeRequest();
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
