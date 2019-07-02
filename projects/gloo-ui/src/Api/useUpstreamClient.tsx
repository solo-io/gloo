import * as React from 'react';
import { requestReducer, Reducer, RequestAction } from './request-reducer';

import {
  ListUpstreamsRequest,
  ListUpstreamsResponse,
  GetUpstreamRequest,
  CreateUpstreamRequest,
  UpdateUpstreamRequest,
  DeleteUpstreamRequest,
  GetUpstreamResponse,
  CreateUpstreamResponse,
  UpdateUpstreamResponse,
  DeleteUpstreamResponse
} from '../proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb_service';

import { upstreamClient } from './grpc-web-hooks';

/* -------------------------------------------------------------------------- */
/*                               LIST UPSTREAMS                               */
/* -------------------------------------------------------------------------- */

export const useGetUpstreamsList = (
  request: ListUpstreamsRequest | null,
  initialData: ListUpstreamsResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<ListUpstreamsResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: ListUpstreamsRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    upstreamClient.listUpstreams(
      request,
      (
        error: ServiceError | null,
        responseMessage: ListUpstreamsResponse | null
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
/*                                GET UPSTREAM                                */
/* -------------------------------------------------------------------------- */

export const useGetUpstream = (
  request: GetUpstreamRequest | null,
  initialData: GetUpstreamResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<GetUpstreamResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: GetUpstreamRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    upstreamClient.getUpstream(
      request,
      (
        error: ServiceError | null,
        responseMessage: GetUpstreamResponse | null
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
/*                               CREATE UPSTREAM                              */
/* -------------------------------------------------------------------------- */

export const useCreateUpstream = (
  request: CreateUpstreamRequest | null,
  initialData: CreateUpstreamResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<CreateUpstreamResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: CreateUpstreamRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    upstreamClient.createUpstream(
      request,
      (
        error: ServiceError | null,
        responseMessage: CreateUpstreamResponse | null
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
/*                               UPDATE UPSTREAM                              */
/* -------------------------------------------------------------------------- */

export const useUpdateUpstream = (
  request: UpdateUpstreamRequest | null,
  initialData: UpdateUpstreamResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<UpdateUpstreamResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: UpdateUpstreamRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    upstreamClient.updateUpstream(
      request,
      (
        error: ServiceError | null,
        responseMessage: UpdateUpstreamResponse | null
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
/*                               DELETE UPSTREAM                              */
/* -------------------------------------------------------------------------- */

export const useDeleteUpstream = (
  request: DeleteUpstreamRequest | null,
  initialData: DeleteUpstreamResponse.AsObject | null = null
) => {
  const [state, dispatch] = React.useReducer<
    Reducer<DeleteUpstreamResponse.AsObject | null>
  >(requestReducer, {
    isLoading: true,
    data: initialData
  });

  const mounted = React.useRef(true);

  const makeRequest = (request: DeleteUpstreamRequest | null) => {
    if (!request) {
      return;
    }
    dispatch({ type: RequestAction.START, payload: null });

    upstreamClient.deleteUpstream(
      request,
      (
        error: ServiceError | null,
        responseMessage: DeleteUpstreamResponse | null
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
