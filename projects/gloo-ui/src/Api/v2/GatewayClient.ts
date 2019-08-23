import { GatewayApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb_service';
import { host } from '../grpc-web-hooks';
import { grpc } from '@improbable-eng/grpc-web';
import {
  ListGatewaysResponse,
  ListGatewaysRequest,
  UpdateGatewayRequest,
  UpdateGatewayResponse,
  GetGatewayResponse,
  GetGatewayRequest
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

export const client = new GatewayApiClient(host, {
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

function getGateway(params: {
  name: string;
  namespace: string;
}): Promise<GetGatewayResponse> {
  return new Promise((resolve, reject) => {
    let req = new GetGatewayRequest();
    req.setRef(getResourceRef(params.name, params.namespace));

    client.getGateway(req, (error, data) => {
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
  delayedCloseTimeout: Duration.AsObject | undefined;
  drainTimeout: Duration.AsObject | undefined;
  generateRequestId: BoolValue.AsObject | undefined;
  idleTimeout: Duration.AsObject | undefined;
  maxRequestHeadersKb: UInt32Value.AsObject | undefined;
  proxy100Continue: boolean;
  requestHeadersForTagsList: string[];
  requestTimeout: Duration.AsObject | undefined;
  serverName: string;
  skipXffAppend: boolean;
  streamIdleTimeout: Duration.AsObject | undefined;
  verbose: boolean;
  useRemoteAddress: BoolValue.AsObject | undefined;
  via: string;
  xffNumTrustedHops: number;
}

function updateGateway(
  params: {
    name: string;
    namespace: string;
    updates: UpdateGatewayHttpData;
  } | null
): Promise<UpdateGatewayResponse> {
  return new Promise((resolve, reject) => {
    let req = new UpdateGatewayRequest();

    if (params !== null) {
      let getReq = new GetGatewayRequest();
      getReq.setRef(getResourceRef(params.name, params.namespace));

      client.getGateway(getReq, (error, data) => {
        if (
          !!data &&
          !!data.getGatewayDetails() &&
          !!data.getGatewayDetails()!.getGateway()
        ) {
          // TODO ~~ This does not work as a final solution
          let updatedGateway = data!.getGatewayDetails()!.getGateway()!;
          const updates = params.updates;

          // If HTTP -- only option currently
          let httpConnectionManagerSettings = new HttpConnectionManagerSettings();

          httpConnectionManagerSettings.setAcceptHttp10(updates.acceptHttp10);
          httpConnectionManagerSettings.setDefaultHostForHttp10(
            updates.defaultHostForHttp10
          );
          setDuration(
            httpConnectionManagerSettings.setDelayedCloseTimeout.bind(
              httpConnectionManagerSettings
            ),
            updates.delayedCloseTimeout
          );
          setDuration(
            httpConnectionManagerSettings.setDrainTimeout.bind(
              httpConnectionManagerSettings
            ),
            updates.drainTimeout
          );
          setBoolVal(
            httpConnectionManagerSettings.setGenerateRequestId.bind(
              httpConnectionManagerSettings
            ),
            updates.generateRequestId
          );
          setDuration(
            httpConnectionManagerSettings.setIdleTimeout.bind(
              httpConnectionManagerSettings
            ),
            updates.idleTimeout
          );
          setUInt32Val(
            httpConnectionManagerSettings.setMaxRequestHeadersKb.bind(
              httpConnectionManagerSettings
            ),
            updates.maxRequestHeadersKb
          );
          if (!!updates.proxy100Continue) {
            httpConnectionManagerSettings.setProxy100Continue(
              updates.proxy100Continue
            );
          }
          setDuration(
            httpConnectionManagerSettings.setRequestTimeout.bind(
              httpConnectionManagerSettings
            ),
            updates.requestTimeout
          );
          httpConnectionManagerSettings.setServerName(updates.serverName);
          httpConnectionManagerSettings.setSkipXffAppend(updates.skipXffAppend);
          setDuration(
            httpConnectionManagerSettings.setStreamIdleTimeout.bind(
              httpConnectionManagerSettings
            ),
            updates.streamIdleTimeout
          );

          let listenerTracing = new ListenerTracingSettings();
          listenerTracing.setRequestHeadersForTagsList(
            updates.requestHeadersForTagsList
          );
          listenerTracing.setVerbose(updates.verbose);
          httpConnectionManagerSettings.setTracing(listenerTracing);

          setBoolVal(
            httpConnectionManagerSettings.setUseRemoteAddress.bind(
              httpConnectionManagerSettings
            ),
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

          updatedGateway.setHttpGateway(httpGateway);

          /*} else if (!!gatewayObj.tcpGateway) {
        // Not visualizing this yet
        // gateway.setTcpGateway(gatewayObj.tcpGateway);
      }*/

          req.setGateway(updatedGateway);

          client.updateGateway(req, (error, data) => {
            console.log({ params, req: req.toObject(), error, data });
            if (error !== null) {
              console.error('Error:', error.message);
              console.error('Code:', error.code);
              console.error('Metadata:', error.metadata);
              reject(error);
            } else {
              resolve(data!);
            }
          });
        }
      });
    } else {
      console.log('empty?');
      reject('null data given');
    }
  });
}

export const gateways = {
  getGatewaysList,
  updateGateway
};
