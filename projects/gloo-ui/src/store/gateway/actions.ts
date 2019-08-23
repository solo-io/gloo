import {
  ListGatewaysRequest,
  ListGatewaysResponse,
  UpdateGatewayRequest,
  UpdateGatewayResponse,
  GetGatewayRequest,
  GetGatewayResponse,
  GatewayDetails
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import { Dispatch } from 'redux';
import { showLoading, hideLoading } from 'react-redux-loading-bar';
import { client } from 'Api/v2/GatewayClient';
import {
  ListGatewaysAction,
  GatewayAction,
  UpdateGatewayAction
} from './types';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { HttpGateway } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v2/gateway_pb';
import { ListenerPlugins } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import {
  BoolValue,
  UInt32Value
} from 'google-protobuf/google/protobuf/wrappers_pb';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { ListenerTracingSettings } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/tracing/tracing_pb';

export function getListGateways(
  listGatewaysRequest: ListGatewaysRequest.AsObject
): Promise<ListGatewaysResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new ListGatewaysRequest();
    request.setNamespacesList(listGatewaysRequest.namespacesList);
    client.listGateways(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        // TODO: normalize
        resolve(data!.toObject());
      }
    });
  });
}

export function getGateway(
  getGatewayRequest: GetGatewayRequest.AsObject
): Promise<GetGatewayResponse> {
  return new Promise((resolve, reject) => {
    let request = new GetGatewayRequest();
    let ref = new ResourceRef();
    ref.setName(getGatewayRequest.ref!.name);
    ref.setNamespace(getGatewayRequest.ref!.namespace);
    request.setRef(ref);
    client.getGateway(request, (error, data) => {
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

export function getUpdateGateway(
  updateGatewayRequest: UpdateGatewayRequest.AsObject
): Promise<UpdateGatewayResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let request = new UpdateGatewayRequest();
    let oldGatewayDetails = await getGateway({
      ref: {
        name: updateGatewayRequest.gateway!.metadata!.name,
        namespace: updateGatewayRequest.gateway!.metadata!.namespace
      }
    });
    let oldGateway = oldGatewayDetails.getGatewayDetails()!.getGateway();
    if (updateGatewayRequest.gateway!.bindAddress) {
      oldGateway!.setBindAddress(updateGatewayRequest.gateway!.bindAddress);
    }
    if (updateGatewayRequest.gateway!.bindPort) {
      oldGateway!.setBindPort(updateGatewayRequest.gateway!.bindPort);
    }

    // TODO: is this changeable?
    if (updateGatewayRequest.gateway!.gatewayProxyName) {
      oldGateway!.setGatewayProxyName(
        updateGatewayRequest.gateway!.gatewayProxyName
      );
    }
    // TODO; merging strategy
    if (updateGatewayRequest.gateway!.httpGateway) {
      let httpGateway = new HttpGateway();
      // TODO HttpListenerPlugins
      if (updateGatewayRequest.gateway!.httpGateway.plugins) {
        let oldHttpListenerPlugins = oldGateway!
          .getHttpGateway()!
          .getPlugins()!;
        let oldHttpConnectionManagerSettings = oldHttpListenerPlugins.getHttpConnectionManagerSettings();
        if (
          updateGatewayRequest.gateway!.httpGateway.plugins
            .httpConnectionManagerSettings
        ) {
          let {
            skipXffAppend,
            via,
            xffNumTrustedHops,
            useRemoteAddress,
            generateRequestId,
            proxy100Continue,
            streamIdleTimeout,
            idleTimeout,
            maxRequestHeadersKb,
            requestTimeout,
            drainTimeout,
            delayedCloseTimeout,
            serverName,
            acceptHttp10,
            defaultHostForHttp10,
            tracing
          } = updateGatewayRequest.gateway!.httpGateway.plugins.httpConnectionManagerSettings!;
          if (skipXffAppend) {
            oldHttpConnectionManagerSettings!.setSkipXffAppend(skipXffAppend);
          }

          if (via) {
            oldHttpConnectionManagerSettings!.setVia(via);
          }

          if (useRemoteAddress) {
            let boolVal = new BoolValue();
            boolVal.setValue(useRemoteAddress.value);
            oldHttpConnectionManagerSettings!.setUseRemoteAddress(boolVal);
          }

          if (xffNumTrustedHops) {
            oldHttpConnectionManagerSettings!.setXffNumTrustedHops(
              xffNumTrustedHops
            );
          }

          if (generateRequestId) {
            let boolVal = new BoolValue();
            boolVal.setValue(generateRequestId.value);
            oldHttpConnectionManagerSettings!.setGenerateRequestId(boolVal);
          }

          // will this not set if false?
          if (proxy100Continue) {
            oldHttpConnectionManagerSettings!.setProxy100Continue(
              proxy100Continue
            );
          }

          if (streamIdleTimeout) {
            let newDuration = new Duration();
            newDuration.setSeconds(streamIdleTimeout.seconds);
            newDuration.setNanos(streamIdleTimeout.nanos);

            oldHttpConnectionManagerSettings!.setStreamIdleTimeout(newDuration);
          }
          if (idleTimeout) {
            let newDuration = new Duration();
            newDuration.setSeconds(idleTimeout.seconds);
            newDuration.setNanos(idleTimeout.nanos);

            oldHttpConnectionManagerSettings!.setIdleTimeout(newDuration);
          }
          if (drainTimeout) {
            let newDuration = new Duration();
            newDuration.setSeconds(drainTimeout.seconds);
            newDuration.setNanos(drainTimeout.nanos);

            oldHttpConnectionManagerSettings!.setDrainTimeout(newDuration);
          }
          if (delayedCloseTimeout) {
            let newDuration = new Duration();
            newDuration.setSeconds(delayedCloseTimeout.seconds);
            newDuration.setNanos(delayedCloseTimeout.nanos);

            oldHttpConnectionManagerSettings!.setDelayedCloseTimeout(
              newDuration
            );
          }

          if (maxRequestHeadersKb) {
            let uInt32 = new UInt32Value();
            uInt32.setValue(maxRequestHeadersKb.value);
            oldHttpConnectionManagerSettings!.setMaxRequestHeadersKb(uInt32);
          }

          if (requestTimeout) {
            let newDuration = new Duration();
            newDuration.setSeconds(requestTimeout.seconds);
            newDuration.setNanos(requestTimeout.nanos);

            oldHttpConnectionManagerSettings!.setStreamIdleTimeout(newDuration);
          }
          if (serverName) {
            oldHttpConnectionManagerSettings!.setServerName(serverName);
          }
          if (acceptHttp10) {
            oldHttpConnectionManagerSettings!.setAcceptHttp10(acceptHttp10);
          }
          if (defaultHostForHttp10) {
            oldHttpConnectionManagerSettings!.setDefaultHostForHttp10(
              defaultHostForHttp10
            );
          }
          if (tracing) {
            let newTracing = new ListenerTracingSettings();
            newTracing.setRequestHeadersForTagsList(
              tracing.requestHeadersForTagsList
            );
            newTracing.setVerbose(tracing.verbose);
            oldHttpConnectionManagerSettings!.setTracing(newTracing);
          }
        }
      }

      let vsRefList = updateGatewayRequest.gateway!.httpGateway.virtualServicesList.map(
        vs => {
          let vsRef = new ResourceRef();
          vsRef.setName(vs.name);
          vsRef.setNamespace(vs.namespace);
          return vsRef;
        }
      );

      httpGateway.setVirtualServicesList(vsRefList);

      oldGateway!.setHttpGateway(httpGateway);
    }
    // TODO
    if (updateGatewayRequest.gateway!.plugins) {
      let oldPlugins = oldGateway!.getPlugins();
      // TODO
      oldGateway!.setPlugins(oldPlugins);
    }
    if (updateGatewayRequest.gateway!.ssl) {
      oldGateway!.setSsl(updateGatewayRequest.gateway!.ssl);
    }
    // TODO
    if (updateGatewayRequest.gateway!.tcpGateway) {
      let oldTcpGateway = oldGateway!.getTcpGateway();
      // find out what changed
      oldGateway!.setTcpGateway(oldTcpGateway);
    }
    if (updateGatewayRequest.gateway!.useProxyProto) {
      let newUseProxyProtoVal = new BoolValue();
      newUseProxyProtoVal.setValue(
        updateGatewayRequest.gateway!.useProxyProto.value
      );
      oldGateway!.setUseProxyProto(newUseProxyProtoVal);
    }

    request.setGateway(oldGateway);
    client.updateGateway(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}

export const listGateways = (
  listGatewaysRequest: ListGatewaysRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());

    try {
      const response = await getListGateways(listGatewaysRequest);
      dispatch<ListGatewaysAction>({
        type: GatewayAction.LIST_GATEWAYS,
        payload: response.gatewayDetailsList
      });
      dispatch(hideLoading());
    } catch (error) {}
  };
};

export const updateGateway = (
  updateGatewayRequest: UpdateGatewayRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await getUpdateGateway(updateGatewayRequest);
      dispatch<UpdateGatewayAction>({
        type: GatewayAction.UPDATE_GATEWAY,
        payload: response.gatewayDetails!
      });
    } catch (error) {}
  };
};
