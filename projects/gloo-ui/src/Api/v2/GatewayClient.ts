import { GatewayApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb_service';
import { host } from '../grpc-web-hooks';
import { grpc } from '@improbable-eng/grpc-web';
import {
  ListGatewaysResponse,
  ListGatewaysRequest,
  UpdateGatewayRequest,
  UpdateGatewayResponse
} from '../../proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import {
  Gateway,
  HttpGateway
} from 'proto/github.com/solo-io/gloo/projects/gateway/api/v2/gateway_pb';
import {
  getResourceRef,
  getStatus,
  getDuration,
  getBoolVal,
  getUInt32Val,
  setMetadata,
  setStatus,
  setDuration,
  setBoolVal,
  setUInt32Val
} from './helpers';
import {
  HttpListenerPlugins,
  ListenerPlugins
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import { GrpcWeb } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/grpc_web/grpc_web_pb';
import { HttpConnectionManagerSettings } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/hcm/hcm_pb';
import {
  BoolValue,
  UInt32Value
} from 'google-protobuf/google/protobuf/wrappers_pb';
import { ListenerTracingSettings } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/tracing/tracing_pb';
import {
  AccessLoggingService,
  AccessLog,
  FileSink
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/als/als_pb';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';

const client = new GatewayApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getGatewaysList(params: {
  namespaces: string[];
}): Promise<ListGatewaysResponse> {
  return new Promise((resolve, reject) => {
    let req = new ListGatewaysRequest();
    req.setNamespacesList(params.namespaces);
    client.listGateways(req, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!);
      }
    });
  });
}

export interface UpdateGatewayHttpData {
  acceptHttp10: boolean;
  defaultHostForHttp10: string;
  delayedCloseTimeout: Duration.AsObject;
  drainTimeout: Duration.AsObject;
  generateRequestId: BoolValue.AsObject;
  idleTimeout: Duration.AsObject;
  maxRequestHeadersKb: UInt32Value.AsObject;
  proxy100Continue: boolean;
  requestHeadersForTagsList: string[];
  requestTimeout: Duration.AsObject;
  serverName: string;
  skipXffAppend: boolean;
  streamIdleTimeout: Duration.AsObject;
  verbose: boolean;
  useRemoteAddress: BoolValue.AsObject;
  via: string;
  xffNumTrustedHops: number;
}

function updateGateway(
  params: {
    originalGateway: Gateway;
    updates: UpdateGatewayHttpData;
  } | null
): Promise<UpdateGatewayResponse> {
  return new Promise((resolve, reject) => {
    let req = new UpdateGatewayRequest();

    if (params !== null) {
      // TODO ~~ This does not work as a final solution
      let updatedGateway = params.originalGateway;
      const updates = params.updates;

      // If HTTP -- only option currently
      let httpConnectionManagerSettings = new HttpConnectionManagerSettings();

      httpConnectionManagerSettings.setAcceptHttp10(updates.acceptHttp10);
      httpConnectionManagerSettings.setDefaultHostForHttp10(
        updates.defaultHostForHttp10
      );
      setDuration(
        httpConnectionManagerSettings.setDelayedCloseTimeout,
        updates.delayedCloseTimeout
      );
      setDuration(
        httpConnectionManagerSettings.setDrainTimeout,
        updates.drainTimeout
      );
      setBoolVal(
        httpConnectionManagerSettings.setGenerateRequestId,
        updates.generateRequestId
      );
      setDuration(
        httpConnectionManagerSettings.setIdleTimeout,
        updates.idleTimeout
      );
      setUInt32Val(
        httpConnectionManagerSettings.setMaxRequestHeadersKb,
        updates.maxRequestHeadersKb
      );
      httpConnectionManagerSettings.setProxy100Continue(
        updates.proxy100Continue
      );
      setDuration(
        httpConnectionManagerSettings.setRequestTimeout,
        updates.requestTimeout
      );
      httpConnectionManagerSettings.setServerName(updates.serverName);
      httpConnectionManagerSettings.setSkipXffAppend(updates.skipXffAppend);
      setDuration(
        httpConnectionManagerSettings.setStreamIdleTimeout,
        updates.streamIdleTimeout
      );

      let listenerTracing = new ListenerTracingSettings();
      listenerTracing.setRequestHeadersForTagsList(
        updates.requestHeadersForTagsList
      );
      listenerTracing.setVerbose(updates.verbose);
      httpConnectionManagerSettings.setTracing(listenerTracing);

      setBoolVal(
        httpConnectionManagerSettings.setUseRemoteAddress,
        updates.useRemoteAddress
      );
      httpConnectionManagerSettings.setVia(updates.via);
      httpConnectionManagerSettings.setXffNumTrustedHops(
        updates.xffNumTrustedHops
      );

      let httpPlugin = new HttpListenerPlugins();
      httpPlugin.setHttpConnectionManagerSettings(
        httpConnectionManagerSettings
      );
      let httpGateway = new HttpGateway();
      httpGateway.setPlugins(httpPlugin);

      /*} else if (!!gatewayObj.tcpGateway) {
        // Not visualizing this yet
        // gateway.setTcpGateway(gatewayObj.tcpGateway);
      }*/

      req.setGateway(updatedGateway);
    }

    client.updateGateway(req, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!);
      }
    });
  });
}

export const gateways = {
  getGatewaysList,
  updateGateway
};
