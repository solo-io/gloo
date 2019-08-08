/* eslint-disable */
import * as React from 'react';

import { grpc } from '@improbable-eng/grpc-web';
import { VirtualServiceApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb_service';
import { UpstreamApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb_service';
import { ConfigApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config_pb_service';
import { SecretApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb_service';

interface ApiResponseType<T> {
  //error?: ServiceError;
  loading: boolean;
  refetch: () => void;
  data: T;
}

export const host = `${
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

export const configClient = new ConfigApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

export const secretClient = new SecretApiClient(host, {
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
