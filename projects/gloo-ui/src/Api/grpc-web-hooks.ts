/* eslint-disable */
import * as React from 'react';

import { grpc } from '@improbable-eng/grpc-web';
import { requestReducer, Reducer, RequestAction } from './request-reducer';
import {
  VirtualServiceApiClient,
  ServiceError
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb_service';
import { ListVirtualServicesRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { ListVirtualServicesResponse } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import {
  ListUpstreamsRequest,
  ListUpstreamsResponse
} from '../proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { UpstreamApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb_service';

interface ApiResponseType<T> {
  //error?: ServiceError;
  loading: boolean;
  refetch: () => void;
  data: T;
}

const host = `${
  process.env.NODE_ENV === 'production'
    ? window.location.origin
    : 'http://localhost:8080'
}`;

export const client = null; /*new GlooEApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});*/

export const virtualServiceClient = new VirtualServiceApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

export const upstreamClient = new UpstreamApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

interface GlooEContextType {
  client: typeof client;
  // other shared bits...
}

export const initialGlooEContext: GlooEContextType = {
  client
};

export const GlooEContext = React.createContext<GlooEContextType>(
  initialGlooEContext
);

export function useGlooEContext() {
  const context = React.useContext(GlooEContext);
  return context;
}

// List Upstreams
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

  const makeRequest = () => {
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

/*
export const useGetNamespacesForMeshes = (
  request: ListNamespacesRequest | null,
  initialData: ListNamespacesResponse | null = null
): ApiResponseType<ListNamespacesResponse> => {
  const [state, dispatch] = React.useReducer<
    Reducer<ListNamespacesResponse | null>
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

    client.listNamespaces(
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
            payload: response
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
}; */
