import { grpc } from '@improbable-eng/grpc-web';
import { Modal } from 'antd';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import {
  BoolValue,
  UInt32Value
} from 'google-protobuf/google/protobuf/wrappers_pb';
import { HttpConnectionManagerSettings } from 'proto/gloo/projects/gloo/api/v1/options/hcm/hcm_pb';
import { ListenerTracingSettings } from 'proto/gloo/projects/gloo/api/v1/options/tracing/tracing_pb';
import { ResourceRef } from 'proto/solo-kit/api/v1/ref_pb';
import {
  GetGatewayRequest,
  GetGatewayResponse,
  ListGatewaysRequest,
  ListGatewaysResponse,
  UpdateGatewayRequest,
  UpdateGatewayResponse,
  UpdateGatewayYamlRequest,
  GatewayDetails
} from 'proto/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import { GatewayApiClient } from 'proto/solo-projects/projects/grpcserver/api/v1/gateway_pb_service';
import { EditedResourceYaml } from 'proto/solo-projects/projects/grpcserver/api/v1/types_pb';
import { host } from 'store';
import { guardByLicense } from 'store/config/actions';
import {
  Gateway,
  HttpGateway,
  TcpGateway
} from 'proto/gloo/projects/gateway/api/v1/gateway_pb';
import { Metadata } from 'proto/solo-kit/api/v1/metadata_pb';
import {
  HttpListenerOptions,
  ListenerOptions
} from 'proto/gloo/projects/gloo/api/v1/options_pb';
import { GrpcWeb } from 'proto/gloo/projects/gloo/api/v1/options/grpc_web/grpc_web_pb';
import { HealthCheck } from 'proto/gloo/projects/gloo/api/v1/options/healthcheck/healthcheck_pb';
import {
  Settings as WafSettings,
  CoreRuleSet
} from 'proto/gloo/projects/gloo/api/v1/enterprise/options/waf/waf_pb';
import {
  FilterConfig,
  DlpRule
} from 'proto/gloo/projects/gloo/api/v1/enterprise/options/dlp/dlp_pb';
import { RuleSet } from 'proto/gloo/projects/gloo/api/external/envoy/extensions/waf/waf_pb';
import { TcpHost } from 'proto/gloo/projects/gloo/api/v1/proxy_pb';
const { warning } = Modal;

const client = new GatewayApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getGateway(
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

function listGateways(): Promise<GatewayDetails.AsObject[]> {
  return new Promise((resolve, reject) => {
    let request = new ListGatewaysRequest();
    client.listGateways(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data?.toObject().gatewayDetailsList);
      }
    });
  });
}

function setGatewayValues(
  gateway: Gateway.AsObject,
  gatewayToUpdate = new Gateway()
): Gateway {
  let {
    metadata,
    ssl,
    bindAddress,
    bindPort,
    status,
    useProxyProto,
    httpGateway,
    tcpGateway,
    proxyNamesList,
    options
  } = gateway;

  if (metadata !== undefined) {
    let { name, namespace } = metadata;
    let newMetadata = new Metadata();
    newMetadata.setName(name);
    newMetadata.setNamespace(namespace);
    gatewayToUpdate.setMetadata(newMetadata);
  }

  if (ssl !== undefined) {
    gatewayToUpdate.setSsl(ssl);
  }

  if (bindAddress !== undefined) {
    gatewayToUpdate.setBindAddress(bindAddress);
  }

  if (bindPort !== undefined) {
    gatewayToUpdate.setBindPort(bindPort);
  }

  if (useProxyProto !== undefined) {
    let useProxyProtoBoolVal = new BoolValue();
    useProxyProtoBoolVal.setValue(useProxyProto.value);
    gatewayToUpdate.setUseProxyProto(useProxyProtoBoolVal);
  }

  if (httpGateway !== undefined) {
    let newHttpGateway = new HttpGateway();
    let {
      virtualServicesList,
      virtualServiceSelectorMap,
      options
    } = httpGateway;
    if (virtualServicesList !== undefined) {
      let vsRefList = virtualServicesList.map(vsRef => {
        let newVsRef = new ResourceRef();
        newVsRef.setName(vsRef.name);
        newVsRef.setNamespace(vsRef.namespace);
        return newVsRef;
      });
      newHttpGateway.setVirtualServicesList(vsRefList);
    }
    // TODO
    if (virtualServiceSelectorMap !== undefined) {
      newHttpGateway.getVirtualServiceSelectorMap();
    }

    if (options !== undefined) {
      let newHttpListenerOptions = new HttpListenerOptions();
      let {
        grpcWeb,
        httpConnectionManagerSettings,
        healthCheck,
        extensions,
        waf,
        dlp
      } = options;

      if (grpcWeb !== undefined) {
        let newGrpcWeb = new GrpcWeb();
        newGrpcWeb.setDisable(grpcWeb?.disable);
        newHttpListenerOptions.setGrpcWeb(newGrpcWeb);
      }

      if (httpConnectionManagerSettings !== undefined) {
        let newHttpConnectionManagerSettings = new HttpConnectionManagerSettings();
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
          tracing,
          forwardClientCertDetails,
          setCurrentClientCertDetails
        } = httpConnectionManagerSettings;
        if (forwardClientCertDetails !== undefined) {
          newHttpConnectionManagerSettings.setForwardClientCertDetails(
            forwardClientCertDetails
          );
        }
        if (setCurrentClientCertDetails !== undefined) {
          let newSetCurrentClientCertDetails = new HttpConnectionManagerSettings.SetCurrentClientCertDetails();
          let { subject, cert, chain, dns, uri } = setCurrentClientCertDetails;
          if (subject !== undefined) {
            let subjectBoolValue = new BoolValue();
            subjectBoolValue.setValue(subject?.value);
            newSetCurrentClientCertDetails.setSubject(subjectBoolValue);
          }
          if (cert !== undefined) {
            newSetCurrentClientCertDetails.setCert(cert);
          }
          if (chain !== undefined) {
            newSetCurrentClientCertDetails.setChain(chain);
          }
          if (dns !== undefined) {
            newSetCurrentClientCertDetails.setDns(dns);
          }
          if (uri !== undefined) {
            newSetCurrentClientCertDetails.setUri(uri);
          }
          newHttpConnectionManagerSettings.setSetCurrentClientCertDetails(
            newSetCurrentClientCertDetails
          );
        }

        if (skipXffAppend !== undefined) {
          newHttpConnectionManagerSettings.setSkipXffAppend(skipXffAppend);
        }

        if (via !== undefined) {
          newHttpConnectionManagerSettings.setVia(via);
        }

        if (useRemoteAddress !== undefined) {
          let boolVal = new BoolValue();
          boolVal.setValue(useRemoteAddress.value);
          newHttpConnectionManagerSettings.setUseRemoteAddress(boolVal);
        }

        if (xffNumTrustedHops !== undefined) {
          newHttpConnectionManagerSettings.setXffNumTrustedHops(
            xffNumTrustedHops
          );
        }

        if (generateRequestId !== undefined) {
          let boolVal = new BoolValue();
          boolVal.setValue(generateRequestId.value);
          newHttpConnectionManagerSettings.setGenerateRequestId(boolVal);
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

          newHttpConnectionManagerSettings.setStreamIdleTimeout(newDuration);
        }
        if (idleTimeout !== undefined) {
          let newDuration = new Duration();
          newDuration.setSeconds(idleTimeout.seconds);
          newDuration.setNanos(idleTimeout.nanos);

          newHttpConnectionManagerSettings.setIdleTimeout(newDuration);
        }
        if (drainTimeout !== undefined) {
          let newDuration = new Duration();
          newDuration.setSeconds(drainTimeout.seconds);
          newDuration.setNanos(drainTimeout.nanos);

          newHttpConnectionManagerSettings.setDrainTimeout(newDuration);
        }
        if (delayedCloseTimeout !== undefined) {
          let newDuration = new Duration();
          newDuration.setSeconds(delayedCloseTimeout.seconds);
          newDuration.setNanos(delayedCloseTimeout.nanos);

          newHttpConnectionManagerSettings.setDelayedCloseTimeout(newDuration);
        }

        if (maxRequestHeadersKb !== undefined) {
          let uInt32 = new UInt32Value();
          uInt32.setValue(maxRequestHeadersKb.value);
          newHttpConnectionManagerSettings.setMaxRequestHeadersKb(uInt32);
        }

        if (requestTimeout !== undefined) {
          let newDuration = new Duration();
          newDuration.setSeconds(requestTimeout.seconds);
          newDuration.setNanos(requestTimeout.nanos);

          newHttpConnectionManagerSettings.setRequestTimeout(newDuration);
        }
        if (serverName !== undefined) {
          newHttpConnectionManagerSettings.setServerName(serverName);
        }
        if (acceptHttp10 !== undefined) {
          newHttpConnectionManagerSettings.setAcceptHttp10(acceptHttp10);
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

      if (healthCheck !== undefined) {
        let newHealthCheck = new HealthCheck();
        newHealthCheck.setPath(healthCheck?.path);
      }

      // TODO: Struct problem,
      if (extensions !== undefined) {
      }

      if (waf !== undefined) {
        let newWafSettings = new WafSettings();
        let {
          disabled,
          customInterventionMessage,
          coreRuleSet,
          ruleSetsList
        } = waf;
        if (disabled !== undefined) {
          newWafSettings.setDisabled(disabled);
        }
        if (customInterventionMessage !== undefined) {
          newWafSettings.setCustomInterventionMessage(
            customInterventionMessage
          );
        }
        if (coreRuleSet !== undefined) {
          let newCoreRuleSet = new CoreRuleSet();
          let { customSettingsString, customSettingsFile } = coreRuleSet;
          if (customSettingsFile !== undefined) {
            newCoreRuleSet.setCustomSettingsFile(customSettingsFile);
          }
          if (customSettingsString !== undefined) {
            newCoreRuleSet.setCustomSettingsString(customSettingsString);
          }
          newWafSettings.setCoreRuleSet(newCoreRuleSet);
        }

        if (ruleSetsList !== undefined) {
          let newRuleSetList = ruleSetsList.map(ruleSet => {
            let newRuleSet = new RuleSet();
            newRuleSet.setRuleStr(ruleSet?.ruleStr);
            newRuleSet.setFilesList(ruleSet?.filesList);
            return newRuleSet;
          });
          newWafSettings.setRuleSetsList(newRuleSetList);
        }
      }
      if (dlp !== undefined) {
        let newDlpFilterConfig = new FilterConfig();
        let { dlpRulesList } = dlp;
        if (dlpRulesList !== undefined) {
          let newDlpRulesList = dlpRulesList.map(dlpRule => {
            let newDlpRule = new DlpRule();
            let { actionsList, matcher } = dlpRule;
            //TODO
            if (matcher !== undefined) {
            }
            // TODO
            if (actionsList !== undefined) {
            }
            return newDlpRule;
          });
          newDlpFilterConfig.setDlpRulesList(newDlpRulesList);
        }
      }
    }
  }

  if (tcpGateway !== undefined) {
    let newTcpGateway = new TcpGateway();
    let { tcpHostsList, options } = tcpGateway;

    if (tcpHostsList !== undefined) {
      let newTcpHostsList = tcpHostsList.map(tcpHost => {
        let newTcpHost = new TcpHost();
        let { name, destination, sslConfig } = tcpHost;
        if (name !== undefined) {
          newTcpHost.setName(name);
        }
        // TODO
        if (destination !== undefined) {
        }
        // TODO
        if (sslConfig !== undefined) {
        }
        return newTcpHost;
      });
      newTcpGateway.setTcpHostsList(newTcpHostsList);
    }

    // TODO
    if (options !== undefined) {
    }
  }

  if (proxyNamesList !== undefined) {
    gatewayToUpdate.setProxyNamesList(proxyNamesList);
  }

  //TODO
  if (options !== undefined) {
    let newListenerOptions = new ListenerOptions();
    let { accessLoggingService, extensions } = options;
  }

  return gatewayToUpdate;
}

function updateGateway(
  updateGatewayRequest: UpdateGatewayRequest.AsObject
): Promise<UpdateGatewayResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let request = new UpdateGatewayRequest();
    let { gateway } = updateGatewayRequest;
    if (gateway !== undefined && gateway.metadata !== undefined) {
      let { name, namespace } = gateway.metadata;
      let gatewayToUpdate = await getGateway({ ref: { name, namespace } });
      let updatedGateway = setGatewayValues(
        gateway,
        gatewayToUpdate.getGatewayDetails()?.getGateway()
      );
      request.setGateway(updatedGateway);
    }
    guardByLicense();
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

function getUpdateGatewayYaml(
  updateGatewayYamlRequest: UpdateGatewayYamlRequest.AsObject
): Promise<UpdateGatewayResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let request = new UpdateGatewayYamlRequest();
    let editedResourceYaml = new EditedResourceYaml();
    let gatewayRef = new ResourceRef();
    gatewayRef.setName(updateGatewayYamlRequest.editedYamlData!.ref!.name);
    gatewayRef.setNamespace(
      updateGatewayYamlRequest.editedYamlData!.ref!.namespace
    );
    editedResourceYaml.setRef(gatewayRef);
    editedResourceYaml.setEditedYaml(
      updateGatewayYamlRequest.editedYamlData!.editedYaml
    );
    request.setEditedYamlData(editedResourceYaml);

    guardByLicense();
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

export const gatewayAPI = {
  listGateways,
  getGateway,
  getUpdateGatewayYaml,
  updateGateway
};
