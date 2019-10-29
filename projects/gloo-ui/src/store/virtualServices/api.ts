import { StringValue } from 'google-protobuf/google/protobuf/wrappers_pb';
import { Route } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/virtual_service_pb';
import { OAuth } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/plugins/extauth/v1/extauth_pb';
import {
  IngressRateLimit,
  RateLimit
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/plugins/ratelimit/ratelimit_pb';
import { DestinationSpec as AwsDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { DestinationSpec as AzureDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { DestinationSpec as GrpcDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/grpc/grpc_pb';
import { DestinationSpec as RestDestinationSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/rest/rest_pb';
import { Parameters } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/transformation/parameters_pb';
import {
  DestinationSpec,
  RoutePlugins
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import {
  Destination,
  RouteAction
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import {
  HeaderMatcher,
  Matcher,
  QueryParameterMatcher,
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/core/matchers/matchers_pb';
import {
  CallCredentials,
  SDSConfig,
  SslConfig,
  SSLFiles,
  SslParameters
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/ssl_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { EditedResourceYaml } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types_pb';
import {
  CreateRouteRequest,
  CreateRouteResponse,
  CreateVirtualServiceRequest,
  CreateVirtualServiceResponse,
  DeleteRouteRequest,
  DeleteRouteResponse,
  DeleteVirtualServiceRequest,
  DeleteVirtualServiceResponse,
  ExtAuthInput,
  GetVirtualServiceRequest,
  GetVirtualServiceResponse,
  IngressRateLimitValue,
  ListVirtualServicesRequest,
  ListVirtualServicesResponse,
  RepeatedRoutes,
  RepeatedStrings,
  RouteInput,
  ShiftRoutesRequest,
  ShiftRoutesResponse,
  SslConfigValue,
  SwapRoutesRequest,
  SwapRoutesResponse,
  UpdateVirtualServiceRequest,
  UpdateVirtualServiceResponse,
  UpdateVirtualServiceYamlRequest,
  VirtualServiceInputV2
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { guardByLicense } from 'store/config/actions';
import { VirtualServiceApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb_service';
import { grpc } from '@improbable-eng/grpc-web';
import { host } from 'store';
import { setInputRouteValues } from 'store/routeTables/api';

const client = new VirtualServiceApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({
    withCredentials: false
  }),
  debug: true
});

function getListVirtualServices(): Promise<
    ListVirtualServicesResponse.AsObject
    > {
  return new Promise((resolve, reject) => {
    let request = new ListVirtualServicesRequest();
    client.listVirtualServices(request, (error, data) => {
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

function getGetVirtualService(
    getVirtualServiceRequest: GetVirtualServiceRequest.AsObject
): Promise<GetVirtualServiceResponse> {
  return new Promise((resolve, reject) => {
    let request = new GetVirtualServiceRequest();
    let ref = new ResourceRef();
    ref.setName(getVirtualServiceRequest.ref!.name);
    ref.setNamespace(getVirtualServiceRequest.ref!.namespace);
    request.setRef(ref);
    client.getVirtualService(request, (error, data) => {
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

function getDeleteVirtualService(
    deleteVirtualServiceRequest: DeleteVirtualServiceRequest.AsObject
): Promise<DeleteVirtualServiceResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new DeleteVirtualServiceRequest();
    let ref = new ResourceRef();
    ref.setName(deleteVirtualServiceRequest.ref!.name);
    ref.setNamespace(deleteVirtualServiceRequest.ref!.namespace);
    request.setRef(ref);
    guardByLicense();
    client.deleteVirtualService(request, (error, data) => {
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

function getUpdateVirtualServiceYaml(
    updateVirtualServiceYamlRequest: UpdateVirtualServiceYamlRequest.AsObject
): Promise<UpdateVirtualServiceResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new UpdateVirtualServiceYamlRequest();
    let vsRef = new ResourceRef();

    let editedYamlData = new EditedResourceYaml();
    let { ref, editedYaml } = updateVirtualServiceYamlRequest.editedYamlData!;
    vsRef.setName(ref!.name);
    vsRef.setNamespace(ref!.namespace);
    editedYamlData.setRef(vsRef);
    editedYamlData.setEditedYaml(editedYaml);

    request.setEditedYamlData(editedYamlData);
    guardByLicense();
    client.updateVirtualServiceYaml(request, (error, data) => {
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

function getDeleteRoute(
    deleteRouteRequest: DeleteRouteRequest.AsObject
): Promise<DeleteRouteResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new DeleteRouteRequest();
    let vsRef = new ResourceRef();
    vsRef.setName(deleteRouteRequest.virtualServiceRef!.name);
    vsRef.setNamespace(deleteRouteRequest.virtualServiceRef!.namespace);
    request.setVirtualServiceRef(vsRef);
    request.setIndex(deleteRouteRequest.index);
    guardByLicense();
    client.deleteRoute(request, (error, data) => {
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

function getSwapRoutes(
    swapRoutesRequest: SwapRoutesRequest.AsObject
): Promise<SwapRoutesResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new SwapRoutesRequest();
    let vsRef = new ResourceRef();
    vsRef.setName(swapRoutesRequest.virtualServiceRef!.name);
    vsRef.setNamespace(swapRoutesRequest.virtualServiceRef!.namespace);

    request.setVirtualServiceRef(vsRef);
    request.setIndex1(swapRoutesRequest.index1);
    request.setIndex2(swapRoutesRequest.index2);
    guardByLicense();
    client.swapRoutes(request, (error, data) => {
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

function getShiftRoutes(
    shiftRoutesRequest: ShiftRoutesRequest.AsObject
): Promise<ShiftRoutesResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new ShiftRoutesRequest();
    let vsRef = new ResourceRef();
    vsRef.setName(shiftRoutesRequest.virtualServiceRef!.name);
    vsRef.setNamespace(shiftRoutesRequest.virtualServiceRef!.namespace);

    request.setVirtualServiceRef(vsRef);
    request.setToIndex(shiftRoutesRequest.toIndex);
    request.setFromIndex(shiftRoutesRequest.fromIndex);
    guardByLicense();
    client.shiftRoutes(request, (error, data) => {
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

function getVirtualServiceForUpdate(
    virtualServiceRef: ResourceRef.AsObject
): Promise<VirtualServiceInputV2> {
  return new Promise(async (resolve, reject) => {
    // input V2
    let inputVirtualServiceV2 = new VirtualServiceInputV2();
    let inputV2Ref = new ResourceRef();
    let inputV2DisplayName = new StringValue(); // updatable
    let inputV2Domains = new RepeatedStrings(); // updatable
    let inputV2Routes = new RepeatedRoutes(); // updatable
    let inputV2SslConfig = new SslConfigValue(); // updatable
    let inputV2RateLimitConfig = new IngressRateLimitValue(); // updatable
    let inputV2ExtAuthConfig = new ExtAuthInput(); // updatable

    // current virtual service info
    let currentVirtualServiceReq = await getGetVirtualService({
      ref: virtualServiceRef
    });

    let currentVirtualServiceDetails = currentVirtualServiceReq.getVirtualServiceDetails();
    if (currentVirtualServiceDetails !== undefined) {
      let currentVirtualService = currentVirtualServiceDetails.getVirtualService();
      if (currentVirtualService !== undefined) {
        let currentVirtualHost = currentVirtualService.getVirtualHost();
        let currentSslConfig = currentVirtualService.getSslConfig();
        let currentDisplayName = currentVirtualService.getDisplayName();
        let currentMetadata = currentVirtualService.getMetadata()!;

        // ref
        inputV2Ref.setName(currentMetadata.getName());
        inputV2Ref.setNamespace(currentMetadata.getNamespace());
        inputVirtualServiceV2.setRef(inputV2Ref);

        // display name
        if (currentDisplayName !== undefined) {
          inputV2DisplayName.setValue(currentDisplayName);
          inputVirtualServiceV2.setDisplayName(inputV2DisplayName);
        }

        if (currentVirtualHost !== undefined) {
          // domains
          let currentDomains = currentVirtualHost.getDomainsList();
          if (currentDomains !== undefined) {
            inputV2Domains.setValuesList(currentDomains);
            inputVirtualServiceV2.setDomains(inputV2Domains);
          }
          // routes
          let currentRoutes = currentVirtualHost.getRoutesList();
          if (currentRoutes !== undefined) {
            inputV2Routes.setValuesList(currentRoutes);
            inputVirtualServiceV2.setRoutes(inputV2Routes);
          }
          // virtual host plugins TODO ?
          let currentVirtualHostPlugins = currentVirtualHost.getVirtualHostPlugins();
          if (currentVirtualHostPlugins !== undefined) {
            let currentVHostExtensions = currentVirtualHostPlugins.getExtensions();
            let currentVHostRetries = currentVirtualHostPlugins.getRetries();
            let currentVHostStats = currentVirtualHostPlugins.getStats();
            let currentVHostHeaderManipulation = currentVirtualHostPlugins.getHeaderManipulation();
            let currentVHostCors = currentVirtualHostPlugins.getCors();
            let currentVHostTransformations = currentVirtualHostPlugins.getTransformations();
          }
        }
        // sslConfig
        if (currentSslConfig !== undefined) {
          inputV2SslConfig.setValue(currentSslConfig);
          inputVirtualServiceV2.setSslConfig(inputV2SslConfig);
        }
      }
      let currentPlugins = currentVirtualServiceDetails.getPlugins();
      if (currentPlugins !== undefined) {
        let currentExtAuth = currentPlugins!.getExtAuth();
        let currentRateLimit = currentPlugins!.getRateLimit();

        if (currentRateLimit !== undefined) {
          inputV2RateLimitConfig.setValue(currentRateLimit.getValue());
        }

        // TODO
        if (currentExtAuth !== undefined) {
          let newExtAuthInputConfig = new ExtAuthInput.Config();
          let currentExtAuthConfig = currentExtAuth!.getValue();
          if (currentExtAuthConfig !== undefined) {
            let newOauth = new OAuth();

            let currentOauth = currentExtAuthConfig.getOauth();
            if (currentOauth !== undefined) {
              newOauth.setClientId(currentOauth.getClientId());
              let clientRef = new ResourceRef();
              // TODO
              newOauth.setClientSecretRef(clientRef);
              newOauth.setIssuerUrl(currentOauth.getIssuerUrl());
              newOauth.setAppUrl(currentOauth.getAppUrl());

              newOauth.setCallbackPath(currentOauth.getCallbackPath());
              newOauth.setScopesList(currentOauth.getScopesList());
            }

            newExtAuthInputConfig.setOauth(newOauth);
          }

          inputV2ExtAuthConfig.setConfig(newExtAuthInputConfig);
          inputVirtualServiceV2.setExtAuthConfig(inputV2ExtAuthConfig);
        }
      }
    } else {
      reject('No current virtual service data');
    }

    resolve(inputVirtualServiceV2);
  });
}

function getUpdateDomains(updateDomainsRequest: {
  ref: ResourceRef.AsObject;
  domains: string[];
}): Promise<UpdateVirtualServiceResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let updateRequest = new UpdateVirtualServiceRequest();

    let inputV2 = await getVirtualServiceForUpdate(updateDomainsRequest.ref);
    let inputV2Domains = new RepeatedStrings();
    if (updateDomainsRequest.domains !== undefined) {
      inputV2Domains.setValuesList(updateDomainsRequest.domains);
    }
    inputV2.setDomains(inputV2Domains);

    updateRequest.setInputV2(inputV2);
    client.updateVirtualService(updateRequest, (error, data) => {
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

function getUpdateDisplayName(updateDisplayNameRequest: {
  ref: ResourceRef.AsObject;
  displayName: string;
}): Promise<UpdateVirtualServiceResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let updateRequest = new UpdateVirtualServiceRequest();

    let inputV2 = await getVirtualServiceForUpdate(
        updateDisplayNameRequest.ref
    );
    let inputV2DisplayName = new StringValue();

    if (updateDisplayNameRequest.displayName !== undefined) {
      inputV2DisplayName.setValue(updateDisplayNameRequest.displayName);
      inputV2.setDisplayName(inputV2DisplayName);
    }
    updateRequest.setInputV2(inputV2);
    client.updateVirtualService(updateRequest, (error, data) => {
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

function getUpdateRoutes(updateRoutesRequest: {
  ref: ResourceRef.AsObject;
  routes: Route.AsObject[];
}): Promise<UpdateVirtualServiceResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let updateRequest = new UpdateVirtualServiceRequest();

    let inputV2 = await getVirtualServiceForUpdate(updateRoutesRequest.ref);
    let inputV2Routes = new RepeatedRoutes();
    let inputRoutesList = updateRoutesRequest.routes!.map(setInputRouteValues);
    inputV2Routes.setValuesList(inputRoutesList);

    inputV2.setRoutes(inputV2Routes);

    updateRequest.setInputV2(inputV2);
    client.updateVirtualService(updateRequest, (error, data) => {
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

function setInputSslConfigValues(sslConfig: SslConfig.AsObject) {
  let inputV2SslConfig = new SslConfigValue();
  if (sslConfig !== undefined) {
    let updatedSslConfig = new SslConfig();
    let {
      secretRef,
      sslFiles,
      sds,
      sniDomainsList,
      verifySubjectAltNameList,
      parameters
    } = sslConfig;
    if (secretRef !== undefined) {
      let sslSecretRef = new ResourceRef();
      sslSecretRef.setName(secretRef!.name);
      sslSecretRef.setNamespace(secretRef!.namespace);
      updatedSslConfig.setSecretRef(sslSecretRef);
    }
    if (sslFiles !== undefined) {
      let updatedSslFiles = new SSLFiles();
      updatedSslFiles.setTlsCert(sslFiles.tlsCert);
      updatedSslFiles.setRootCa(sslFiles.rootCa);
      updatedSslFiles.setTlsKey(sslFiles.tlsKey);
      updatedSslConfig.setSslFiles(updatedSslFiles);
    }
    if (sds !== undefined) {
      let {
        callCredentials,
        certificatesSecretName,
        targetUri,
        validationContextName
      } = sds;
      let updatedSds = new SDSConfig();
      if (callCredentials !== undefined) {
        let newCallCreds = new CallCredentials();
        let newCallCredsSource = new CallCredentials.FileCredentialSource();
        newCallCredsSource.setHeader(
            callCredentials.fileCredentialSource!.header
        );
        newCallCredsSource.setTokenFileName(
            callCredentials.fileCredentialSource!.tokenFileName
        );
        newCallCreds.setFileCredentialSource(newCallCredsSource);
        updatedSds.setCallCredentials(newCallCreds);
      }
      updatedSds.setCertificatesSecretName(certificatesSecretName);
      updatedSds.setTargetUri(targetUri);
      updatedSds.setValidationContextName(validationContextName);
      updatedSslConfig.setSds(updatedSds);
    }
    if (sniDomainsList !== undefined) {
      updatedSslConfig.setSniDomainsList(sniDomainsList);
    }
    if (verifySubjectAltNameList !== undefined) {
      updatedSslConfig.setVerifySubjectAltNameList(verifySubjectAltNameList);
    }
    if (parameters !== undefined) {
      let {
        minimumProtocolVersion,
        maximumProtocolVersion,
        cipherSuitesList,
        ecdhCurvesList
      } = parameters;
      let newSslParams = new SslParameters();

      newSslParams.setCipherSuitesList(cipherSuitesList);
      newSslParams.setEcdhCurvesList(ecdhCurvesList);
      newSslParams.setMaximumProtocolVersion(maximumProtocolVersion);
      newSslParams.setMinimumProtocolVersion(minimumProtocolVersion);
      updatedSslConfig.setParameters(newSslParams);
    }

    inputV2SslConfig.setValue(updatedSslConfig);
  }
  return inputV2SslConfig;
}
function getUpdateSslConfig(updateSslConfigRequest: {
  ref: ResourceRef.AsObject;
  sslConfig: SslConfig.AsObject;
}): Promise<UpdateVirtualServiceResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let updateRequest = new UpdateVirtualServiceRequest();

    let inputV2 = await getVirtualServiceForUpdate(updateSslConfigRequest.ref);
    let inputV2SslConfig = setInputSslConfigValues(
        updateSslConfigRequest.sslConfig!
    );

    inputV2.setSslConfig(inputV2SslConfig);
    updateRequest.setInputV2(inputV2);
    client.updateVirtualService(updateRequest, (error, data) => {
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

function getCreateRoute(
    createRouteRequest: CreateRouteRequest.AsObject
): Promise<CreateRouteResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let createRequest = new CreateRouteRequest();
    let inputRoute = new RouteInput();

    let { virtualServiceRef, index, route } = createRouteRequest.input!;

    if (route !== undefined) {
      let newRoute = setInputRouteValues(route);

      inputRoute.setRoute(newRoute);
    }
    inputRoute.setIndex(index);
    if (virtualServiceRef !== undefined) {
      let vsRef = new ResourceRef();
      vsRef.setName(virtualServiceRef!.name);
      vsRef.setNamespace(virtualServiceRef!.namespace);
      inputRoute.setVirtualServiceRef(vsRef);
    }

    createRequest.setInput(inputRoute);
    client.createRoute(createRequest, (error, data) => {
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

function getUpdateVirtualService(
    updateVirtualServiceRequest: UpdateVirtualServiceRequest.AsObject
): Promise<UpdateVirtualServiceResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let updateRequest = new UpdateVirtualServiceRequest();
    // input V2
    let inputV2 = new VirtualServiceInputV2();
    let inputV2Ref = new ResourceRef();
    let inputV2DisplayName = new StringValue(); //updatable
    let inputV2Domains = new RepeatedStrings(); // updatable
    let inputV2Routes = new RepeatedRoutes(); // updatable
    let inputV2SslConfig = new SslConfigValue(); // updatable
    let inputV2RateLimitConfig = new IngressRateLimitValue(); // updatable
    let inputV2ExtAuthConfig = new ExtAuthInput(); // updatable
    if (updateVirtualServiceRequest.inputV2 !== undefined) {
      let inputSslConfig = updateVirtualServiceRequest.inputV2!.sslConfig;
      let inputDisplayName = updateVirtualServiceRequest.inputV2!.displayName;

      // ref
      inputV2Ref.setName(updateVirtualServiceRequest.inputV2!.ref!.name);
      inputV2Ref.setNamespace(
          updateVirtualServiceRequest.inputV2!.ref!.namespace
      );
      inputV2.setRef(inputV2Ref);
      // display name
      if (inputDisplayName !== undefined) {
        inputV2DisplayName.setValue(inputDisplayName.value);
        inputV2.setDisplayName(inputV2DisplayName);
      }
      // domains
      let inputDomains = updateVirtualServiceRequest.inputV2!.domains;
      if (inputDomains !== undefined) {
        inputV2Domains.setValuesList(inputDomains!.valuesList);
        inputV2.setDomains(inputV2Domains);
      }
      //routes
      let inputRoutes = updateVirtualServiceRequest.inputV2!.routes;
      if (inputRoutes !== undefined) {
        inputV2Routes.setValuesList(
            inputRoutes!.valuesList.map(setInputRouteValues)
        );
        inputV2.setRoutes(inputV2Routes);
      }
      //extAuth
      if (updateVirtualServiceRequest.inputV2!.extAuthConfig !== undefined) {
        inputV2ExtAuthConfig = setInputExtAuthValues(
            updateVirtualServiceRequest.inputV2!.extAuthConfig
        );
        inputV2.setExtAuthConfig(inputV2ExtAuthConfig);
      }

      // rate limit
      if (updateVirtualServiceRequest.inputV2!.rateLimitConfig !== undefined) {
        inputV2RateLimitConfig = setInputRateLimitValues(
            updateVirtualServiceRequest.inputV2!.rateLimitConfig.value!
        );
        inputV2.setRateLimitConfig(inputV2RateLimitConfig);
      }

      // // virtual host plugins TODO ?
      // let inputVirtualHostPlugins = inputVirtualHost.getVirtualHostPlugins();
      // if (inputVirtualHostPlugins !== undefined) {
      //   let currentVHostExtensions = inputVirtualHostPlugins.getExtensions();
      //   let currentVHostRetries = inputVirtualHostPlugins.getRetries();
      //   let currentVHostStats = inputVirtualHostPlugins.getStats();
      //   let currentVHostHeaderManipulation = inputVirtualHostPlugins.getHeaderManipulation();
      //   let currentVHostCors = inputVirtualHostPlugins.getCors();
      //   let currentVHostTransformations = inputVirtualHostPlugins.getTransformations();
      // }

      // sslConfig
      if (inputSslConfig !== undefined && inputSslConfig.value !== undefined) {
        inputV2SslConfig = setInputSslConfigValues(inputSslConfig!.value);
      }
      inputV2.setSslConfig(inputV2SslConfig);
    }
    updateRequest.setInputV2(inputV2);
    client.updateVirtualService(updateRequest, (error, data) => {
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

function getCreateVirtualService(
    createVirtualSeviceRequest: CreateVirtualServiceRequest.AsObject
): Promise<CreateVirtualServiceResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let createRequest = new CreateVirtualServiceRequest();
    // input V2
    let inputV2 = new VirtualServiceInputV2();
    let inputV2Ref = new ResourceRef();
    let inputV2DisplayName = new StringValue(); // updatable
    let inputV2Domains = new RepeatedStrings(); // updatable
    let inputV2Routes = new RepeatedRoutes(); // updatable
    let inputV2SslConfig = new SslConfigValue(); // updatable
    let inputV2RateLimitConfig = new IngressRateLimitValue(); // updatable
    let inputV2ExtAuthConfig = new ExtAuthInput(); // updatable

    if (createVirtualSeviceRequest.inputV2 !== undefined) {
      let {
        extAuthConfig,
        rateLimitConfig,
        ref,
        routes,
        displayName,
        domains,
        sslConfig
      } = createVirtualSeviceRequest.inputV2!;

      // ref
      inputV2Ref.setName(ref!.name);
      inputV2Ref.setNamespace(ref!.namespace);
      inputV2.setRef(inputV2Ref);

      // display name
      if (displayName !== undefined) {
        inputV2DisplayName.setValue(displayName.value);
        inputV2.setDisplayName(inputV2DisplayName);
      }

      // domains
      if (domains !== undefined) {
        inputV2Domains.setValuesList(domains.valuesList);
        inputV2.setDomains(inputV2Domains);
      }

      // virtual host plugins TODO ?
      // let currentVirtualHostPlugins = currentVirtualHost.getVirtualHostPlugins();
      // if (currentVirtualHostPlugins !== undefined) {
      //   let currentVHostExtensions = currentVirtualHostPlugins.getExtensions();
      //   let currentVHostRetries = currentVirtualHostPlugins.getRetries();
      //   let currentVHostStats = currentVirtualHostPlugins.getStats();
      //   let currentVHostHeaderManipulation = currentVirtualHostPlugins.getHeaderManipulation();
      //   let currentVHostCors = currentVirtualHostPlugins.getCors();
      //   let currentVHostTransformations = currentVirtualHostPlugins.getTransformations();
      // }

      //routes
      if (routes !== undefined) {
        inputV2Routes.setValuesList(
            routes!.valuesList.map(setInputRouteValues)
        );
        inputV2.setRoutes(inputV2Routes);
      }
      //extAuth
      if (extAuthConfig !== undefined) {
        inputV2ExtAuthConfig = setInputExtAuthValues(extAuthConfig);
        inputV2.setExtAuthConfig(inputV2ExtAuthConfig);
      }

      // rate limit
      if (rateLimitConfig !== undefined) {
        inputV2RateLimitConfig = setInputRateLimitValues(
            rateLimitConfig.value!
        );
        inputV2.setRateLimitConfig(inputV2RateLimitConfig);
      }
      // sslConfig
      if (sslConfig !== undefined) {
        inputV2SslConfig = setInputSslConfigValues(sslConfig!.value!);
        inputV2.setSslConfig(inputV2SslConfig);
      }
    }

    createRequest.setInputV2(inputV2);
    client.createVirtualService(createRequest, (error, data) => {
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

function setInputExtAuthValues(extAuthConfig: ExtAuthInput.AsObject) {
  let inputV2ExtAuthConfig = new ExtAuthInput(); // updatable

  // ext auth config
  if (extAuthConfig.config !== undefined) {
    let newExtAuthConfig = new ExtAuthInput.Config();
    if (extAuthConfig.config.oauth !== undefined) {
      let newOauthConfig = newExtAuthConfig.getOauth();
      if (newOauthConfig !== undefined) {
        let {
          clientId,
          clientSecretRef,
          issuerUrl,
          appUrl,
          callbackPath,
          scopesList
        } = extAuthConfig.config.oauth;
        if (clientId !== undefined) {
          newOauthConfig.setClientId(clientId);
        }
        if (clientSecretRef !== undefined) {
          let secretRef = new ResourceRef();
          secretRef.setName(clientSecretRef.name);
          secretRef.setNamespace(clientSecretRef.namespace);

          newOauthConfig.setClientSecretRef(secretRef);
        }
        if (issuerUrl !== undefined) {
          newOauthConfig.setIssuerUrl(issuerUrl);
        }
        if (appUrl !== undefined) {
          newOauthConfig.setAppUrl(appUrl);
        }
        if (callbackPath !== undefined) {
          newOauthConfig.setCallbackPath(callbackPath);
        }
        if (scopesList !== undefined) {
          newOauthConfig.setScopesList(scopesList);
        }
      }
      newExtAuthConfig.setOauth(newOauthConfig);
    }
    if (extAuthConfig.config.customAuth !== undefined) {
      // TODO
    }

    inputV2ExtAuthConfig.setConfig(newExtAuthConfig);
  }
  return inputV2ExtAuthConfig;
}

function getUpdateExtAuth(updateExtAuthRequest: {
  ref: ResourceRef.AsObject;
  extAuthConfig: ExtAuthInput.AsObject;
}): Promise<UpdateVirtualServiceResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let updateRequest = new UpdateVirtualServiceRequest();
    // input V2
    let inputV2 = await getVirtualServiceForUpdate(updateExtAuthRequest.ref);

    let inputV2ExtAuthConfig = setInputExtAuthValues(
        updateExtAuthRequest.extAuthConfig
    );
    inputV2.setExtAuthConfig(inputV2ExtAuthConfig);
    updateRequest.setInputV2(inputV2);
    client.updateVirtualService(updateRequest, (error, data) => {
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

function setInputRateLimitValues(rateLimitConfig: IngressRateLimit.AsObject) {
  let inputV2RateLimitConfig = new IngressRateLimitValue();

  // rate limit config
  if (rateLimitConfig !== undefined) {
    let newIngressRateLimit = new IngressRateLimit();
    let { anonymousLimits, authorizedLimits } = rateLimitConfig!;

    if (authorizedLimits !== undefined) {
      let inputAuthorizedRateLimit = new RateLimit();
      inputAuthorizedRateLimit.setRequestsPerUnit(
          authorizedLimits.requestsPerUnit
      );
      inputAuthorizedRateLimit.setUnit(authorizedLimits.unit);
      newIngressRateLimit.setAuthorizedLimits(inputAuthorizedRateLimit);
      inputV2RateLimitConfig.setValue(newIngressRateLimit);
    }
    if (anonymousLimits !== undefined) {
      let inputAnonymousRateLimit = new RateLimit();
      inputAnonymousRateLimit.setRequestsPerUnit(
          anonymousLimits.requestsPerUnit
      );
      inputAnonymousRateLimit.setUnit(anonymousLimits.unit);
      newIngressRateLimit.setAnonymousLimits(inputAnonymousRateLimit);
      inputV2RateLimitConfig.setValue(newIngressRateLimit);
    }
  }
  return inputV2RateLimitConfig;
}
function getUpdateRateLimit(updateRateLimitRequest: {
  ref: ResourceRef.AsObject;
  rateLimitConfig: IngressRateLimit.AsObject;
}): Promise<UpdateVirtualServiceResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let updateRequest = new UpdateVirtualServiceRequest();

    let inputV2 = await getVirtualServiceForUpdate(updateRateLimitRequest.ref);

    let inputV2RateLimitConfig = setInputRateLimitValues(
        updateRateLimitRequest.rateLimitConfig
    );

    inputV2.setRateLimitConfig(inputV2RateLimitConfig);
    updateRequest.setInputV2(inputV2);
    client.updateVirtualService(updateRequest, (error, data) => {
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

export const virtualServices = {
  getListVirtualServices,
  getGetVirtualService,
  getDeleteVirtualService,
  getUpdateVirtualServiceYaml,
  getDeleteRoute,
  getSwapRoutes,
  getShiftRoutes,
  getVirtualServiceForUpdate,
  getUpdateDomains,
  getUpdateDisplayName,
  setInputRouteValues,
  getUpdateRoutes,
  setInputSslConfigValues,
  getUpdateSslConfig,
  getCreateRoute,
  getUpdateVirtualService,
  getCreateVirtualService,
  setInputExtAuthValues,
  getUpdateExtAuth,
  setInputRateLimitValues,
  getUpdateRateLimit
};