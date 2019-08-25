import { client } from 'Api/v2/GatewayClient';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import {
  BoolValue,
  UInt32Value
} from 'google-protobuf/google/protobuf/wrappers_pb';
import { HttpConnectionManagerSettings } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/hcm/hcm_pb';
import { ListenerTracingSettings } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/tracing/tracing_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  GetGatewayRequest,
  GetGatewayResponse,
  ListGatewaysRequest,
  ListGatewaysResponse,
  UpdateGatewayRequest,
  UpdateGatewayResponse
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import { hideLoading, showLoading } from 'react-redux-loading-bar';
import { Dispatch } from 'redux';
import {
  GatewayAction,
  ListGatewaysAction,
  UpdateGatewayAction
} from './types';

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
  console.log('updateGatewayRequest', updateGatewayRequest);
  return new Promise(async (resolve, reject) => {
    let request = new UpdateGatewayRequest();
    let oldGatewayDetails = await getGateway({
      ref: {
        name: updateGatewayRequest.gateway!.metadata!.name,
        namespace: updateGatewayRequest.gateway!.metadata!.namespace
      }
    });
    let oldGatewayD = oldGatewayDetails.getGatewayDetails();
    if (oldGatewayD !== undefined) {
      let oldGateway = oldGatewayD.getGateway();

      // if (updateGatewayRequest.gateway!.bindAddress) {
      //   oldGateway!.setBindAddress(updateGatewayRequest.gateway!.bindAddress);
      // }
      // if (updateGatewayRequest.gateway!.bindPort) {
      //   oldGateway!.setBindPort(updateGatewayRequest.gateway!.bindPort);
      // }

      // TODO: is this changeable?
      // if (updateGatewayRequest.gateway!.gatewayProxyName) {
      //   oldGateway!.setGatewayProxyName(
      //     updateGatewayRequest.gateway!.gatewayProxyName
      //   );
      // }
      // TODO; merging strategy
      if (oldGateway !== undefined) {
        if (oldGateway.getHttpGateway() !== undefined) {
          let oldHttpGateway = oldGateway.getHttpGateway();
          if (
            oldHttpGateway !== undefined &&
            oldHttpGateway.getPlugins() !== undefined
          ) {
            let oldHttpPlugins = oldHttpGateway.getPlugins();
            if (oldHttpPlugins !== undefined) {
              let newHttpConnectionManagerSettings = new HttpConnectionManagerSettings();
              if (
                updateGatewayRequest.gateway !== undefined &&
                updateGatewayRequest.gateway.httpGateway !== undefined
              ) {
                if (
                  updateGatewayRequest.gateway.httpGateway.plugins !== undefined
                ) {
                  if (
                    updateGatewayRequest.gateway.httpGateway.plugins
                      .httpConnectionManagerSettings !== undefined
                  ) {
                    // set new Httpconnectionmaneg
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
                    } = updateGatewayRequest.gateway.httpGateway.plugins.httpConnectionManagerSettings;
                    if (skipXffAppend !== undefined) {
                      newHttpConnectionManagerSettings.setSkipXffAppend(
                        skipXffAppend
                      );
                    }

                    if (via !== undefined) {
                      newHttpConnectionManagerSettings.setVia(via);
                    }

                    if (useRemoteAddress !== undefined) {
                      let boolVal = new BoolValue();
                      boolVal.setValue(useRemoteAddress.value);
                      newHttpConnectionManagerSettings.setUseRemoteAddress(
                        boolVal
                      );
                    }

                    if (xffNumTrustedHops !== undefined) {
                      newHttpConnectionManagerSettings.setXffNumTrustedHops(
                        xffNumTrustedHops
                      );
                    }

                    if (generateRequestId !== undefined) {
                      let boolVal = new BoolValue();
                      boolVal.setValue(generateRequestId.value);
                      newHttpConnectionManagerSettings.setGenerateRequestId(
                        boolVal
                      );
                    }

                    if (proxy100Continue !== undefined) {
                      newHttpConnectionManagerSettings.setProxy100Continue(
                        proxy100Continue
                      );
                    }

                    if (streamIdleTimeout !== undefined) {
                      let newDuration = new Duration();
                      newDuration.setSeconds(streamIdleTimeout.seconds);
                      newDuration.setNanos(streamIdleTimeout.nanos);

                      newHttpConnectionManagerSettings.setStreamIdleTimeout(
                        newDuration
                      );
                    }
                    if (idleTimeout !== undefined) {
                      let newDuration = new Duration();
                      newDuration.setSeconds(idleTimeout.seconds);
                      newDuration.setNanos(idleTimeout.nanos);

                      newHttpConnectionManagerSettings.setIdleTimeout(
                        newDuration
                      );
                    }
                    if (drainTimeout !== undefined) {
                      let newDuration = new Duration();
                      newDuration.setSeconds(drainTimeout.seconds);
                      newDuration.setNanos(drainTimeout.nanos);

                      newHttpConnectionManagerSettings.setDrainTimeout(
                        newDuration
                      );
                    }
                    if (delayedCloseTimeout !== undefined) {
                      let newDuration = new Duration();
                      newDuration.setSeconds(delayedCloseTimeout.seconds);
                      newDuration.setNanos(delayedCloseTimeout.nanos);

                      newHttpConnectionManagerSettings.setDelayedCloseTimeout(
                        newDuration
                      );
                    }

                    if (maxRequestHeadersKb !== undefined) {
                      let uInt32 = new UInt32Value();
                      uInt32.setValue(maxRequestHeadersKb.value);
                      newHttpConnectionManagerSettings.setMaxRequestHeadersKb(
                        uInt32
                      );
                    }

                    if (requestTimeout !== undefined) {
                      let newDuration = new Duration();
                      newDuration.setSeconds(requestTimeout.seconds);
                      newDuration.setNanos(requestTimeout.nanos);

                      newHttpConnectionManagerSettings.setStreamIdleTimeout(
                        newDuration
                      );
                    }
                    if (serverName !== undefined) {
                      newHttpConnectionManagerSettings.setServerName(
                        serverName
                      );
                    }
                    if (acceptHttp10 !== undefined) {
                      newHttpConnectionManagerSettings.setAcceptHttp10(
                        acceptHttp10
                      );
                    }
                    if (defaultHostForHttp10 !== undefined) {
                      newHttpConnectionManagerSettings.setDefaultHostForHttp10(
                        defaultHostForHttp10
                      );
                    }
                    if (tracing !== undefined) {
                      let newTracing = new ListenerTracingSettings();
                      newTracing.setRequestHeadersForTagsList(
                        tracing.requestHeadersForTagsList
                      );
                      newTracing.setVerbose(tracing.verbose);
                      newHttpConnectionManagerSettings.setTracing(newTracing);
                    }
                  }
                }
              }

              oldHttpPlugins.setHttpConnectionManagerSettings(
                newHttpConnectionManagerSettings
              );
            }
            oldHttpGateway.setPlugins(oldHttpPlugins);
            oldGateway.setHttpGateway(oldHttpGateway);
          }
        }

        //   // let vsRefList = updateGatewayRequest.gateway!.httpGateway.virtualServicesList.map(
        //   //   vs => {
        //   //     let vsRef = new ResourceRef();
        //   //     vsRef.setName(vs.name);
        //   //     vsRef.setNamespace(vs.namespace);
        //   //     return vsRef;
        //   //   }
        //   // )

        //   // httpGateway.setVirtualServicesList(vsRefList);

        //   // oldGateway!.setHttpGateway(httpGateway);
        // }
        // TODO
        if (updateGatewayRequest.gateway !== undefined) {
          if (updateGatewayRequest.gateway.plugins) {
            let oldPlugins = oldGateway.getPlugins();
            // TODO
            oldGateway.setPlugins(oldPlugins);
          }
          if (updateGatewayRequest.gateway.ssl) {
            oldGateway.setSsl(updateGatewayRequest.gateway.ssl);
          }
          // TODO
          if (updateGatewayRequest.gateway.tcpGateway) {
            let oldTcpGateway = oldGateway.getTcpGateway();
            // find out what changed
            oldGateway.setTcpGateway(oldTcpGateway);
          }
          if (updateGatewayRequest.gateway.useProxyProto) {
            let newUseProxyProtoVal = new BoolValue();
            newUseProxyProtoVal.setValue(
              updateGatewayRequest.gateway.useProxyProto.value
            );
            oldGateway.setUseProxyProto(newUseProxyProtoVal);
          }
        }

        request.setGateway(oldGateway);
      }
    }
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
