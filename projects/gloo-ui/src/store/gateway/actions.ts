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
  UpdateGatewayResponse,
  UpdateGatewayYamlRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import { hideLoading, showLoading } from 'react-redux-loading-bar';
import { Dispatch } from 'redux';
import {
  GatewayAction,
  ListGatewaysAction,
  UpdateGatewayAction,
  UpdateGatewayYamlAction
} from './types';
import { HttpListenerPlugins } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import { EditedResourceYaml } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types_pb';
import { getResourceRef } from 'Api/v2/helpers';
import { Modal } from 'antd';
import { SuccessMessageAction, MessageAction } from 'store/modal/types';
const { warning } = Modal;

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
      let currentGateway = oldGatewayD.getGateway();

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

      if (currentGateway !== undefined) {
        if (currentGateway.getHttpGateway() !== undefined) {
          let httpGatewayToUpdate = currentGateway.getHttpGateway()!;
          let oldHttpPlugins = httpGatewayToUpdate.getPlugins();
          if (oldHttpPlugins === undefined) {
            oldHttpPlugins = new HttpListenerPlugins();
          }
          let oldHttpCMS = oldHttpPlugins.getHttpConnectionManagerSettings();
          if (oldHttpCMS === undefined) {
            oldHttpCMS = new HttpConnectionManagerSettings();
          }

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
                  oldHttpCMS.setSkipXffAppend(skipXffAppend);
                }

                if (via !== undefined) {
                  oldHttpCMS.setVia(via);
                }

                if (useRemoteAddress !== undefined) {
                  let boolVal = new BoolValue();
                  boolVal.setValue(useRemoteAddress.value);
                  oldHttpCMS.setUseRemoteAddress(boolVal);
                }

                if (xffNumTrustedHops !== undefined) {
                  oldHttpCMS.setXffNumTrustedHops(xffNumTrustedHops);
                }

                if (generateRequestId !== undefined) {
                  let boolVal = new BoolValue();
                  boolVal.setValue(generateRequestId.value);
                  oldHttpCMS.setGenerateRequestId(boolVal);
                }

                if (proxy100Continue !== undefined) {
                  oldHttpCMS.setProxy100Continue(proxy100Continue);
                }

                if (streamIdleTimeout !== undefined) {
                  let newDuration = new Duration();
                  newDuration.setSeconds(streamIdleTimeout.seconds);
                  newDuration.setNanos(streamIdleTimeout.nanos);

                  oldHttpCMS.setStreamIdleTimeout(newDuration);
                }
                if (idleTimeout !== undefined) {
                  let newDuration = new Duration();
                  newDuration.setSeconds(idleTimeout.seconds);
                  newDuration.setNanos(idleTimeout.nanos);

                  oldHttpCMS.setIdleTimeout(newDuration);
                }
                if (drainTimeout !== undefined) {
                  let newDuration = new Duration();
                  newDuration.setSeconds(drainTimeout.seconds);
                  newDuration.setNanos(drainTimeout.nanos);

                  oldHttpCMS.setDrainTimeout(newDuration);
                }
                if (delayedCloseTimeout !== undefined) {
                  let newDuration = new Duration();
                  newDuration.setSeconds(delayedCloseTimeout.seconds);
                  newDuration.setNanos(delayedCloseTimeout.nanos);

                  oldHttpCMS.setDelayedCloseTimeout(newDuration);
                }

                if (maxRequestHeadersKb !== undefined) {
                  let uInt32 = new UInt32Value();
                  uInt32.setValue(maxRequestHeadersKb.value);
                  oldHttpCMS.setMaxRequestHeadersKb(uInt32);
                }

                if (requestTimeout !== undefined) {
                  let newDuration = new Duration();
                  newDuration.setSeconds(requestTimeout.seconds);
                  newDuration.setNanos(requestTimeout.nanos);

                  oldHttpCMS.setRequestTimeout(newDuration);
                }
                if (serverName !== undefined) {
                  oldHttpCMS.setServerName(serverName);
                }
                if (acceptHttp10 !== undefined) {
                  oldHttpCMS.setAcceptHttp10(acceptHttp10);
                }
                if (defaultHostForHttp10 !== undefined) {
                  oldHttpCMS.setDefaultHostForHttp10(defaultHostForHttp10);
                }
                if (tracing !== undefined) {
                  let newTracing = new ListenerTracingSettings();
                  newTracing.setRequestHeadersForTagsList(
                    tracing.requestHeadersForTagsList
                  );
                  newTracing.setVerbose(tracing.verbose);
                  oldHttpCMS.setTracing(newTracing);
                }
              }
            }
          }

          oldHttpPlugins.setHttpConnectionManagerSettings(oldHttpCMS);
          httpGatewayToUpdate.setPlugins(oldHttpPlugins);

          currentGateway.setHttpGateway(httpGatewayToUpdate);
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
            let oldPlugins = currentGateway.getPlugins();
            // TODO
            currentGateway.setPlugins(oldPlugins);
          }
          if (updateGatewayRequest.gateway.ssl) {
            currentGateway.setSsl(updateGatewayRequest.gateway.ssl);
          }
          // TODO
          if (updateGatewayRequest.gateway.tcpGateway) {
            let oldTcpGateway = currentGateway.getTcpGateway();
            // find out what changed
            currentGateway.setTcpGateway(oldTcpGateway);
          }
          if (updateGatewayRequest.gateway.useProxyProto) {
            let newUseProxyProtoVal = new BoolValue();
            newUseProxyProtoVal.setValue(
              updateGatewayRequest.gateway.useProxyProto.value
            );
            currentGateway.setUseProxyProto(newUseProxyProtoVal);
          }
        }

        request.setGateway(currentGateway);
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

export function getUpdateGatewayYaml(
  updateGatewayYamlRequest: UpdateGatewayYamlRequest.AsObject
): Promise<UpdateGatewayResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let request = new UpdateGatewayYamlRequest();
    let editedResourceYaml = new EditedResourceYaml();
    editedResourceYaml.setRef(
      getResourceRef(
        updateGatewayYamlRequest.editedYamlData!.ref!.name,
        updateGatewayYamlRequest.editedYamlData!.ref!.namespace
      )
    );
    editedResourceYaml.setEditedYaml(
      updateGatewayYamlRequest.editedYamlData!.editedYaml
    );
    request.setEditedYamlData(editedResourceYaml);

    client.updateGatewayYaml(request, (error, data) => {
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

      dispatch<SuccessMessageAction>({
        type: MessageAction.SUCCESS_MESSAGE,
        message: 'Gateway successfully updated.'
      });
      dispatch(hideLoading());
    } catch (error) {
      warning({
        title: 'There was an error updating the gateway configuration.',
        content: error.message
      });
    }
  };
};

export const updateGatewayYaml = (
  updateGatewayYamlRequest: UpdateGatewayYamlRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());
    try {
      const response = await getUpdateGatewayYaml(updateGatewayYamlRequest);
      dispatch<UpdateGatewayYamlAction>({
        type: GatewayAction.UPDATE_GATEWAY_YAML,
        payload: response.gatewayDetails!
      });
    } catch (error) {
      warning({
        title: 'There was an error updating the gateway configuration.',
        content: error.message
      });
    }
  };
};
