import { grpc } from '@improbable-eng/grpc-web';
import { AwsValuesType } from 'Components/Features/Upstream/Creation/AwsUpstreamForm';
import { AzureValuesType } from 'Components/Features/Upstream/Creation/AzureUpstreamForm';
import { ConsulVauesType } from 'Components/Features/Upstream/Creation/ConsulUpstreamForm';
import { KubeValuesType } from 'Components/Features/Upstream/Creation/KubeUpstreamForm';
import { StaticValuesType } from 'Components/Features/Upstream/Creation/StaticUpstreamForm';
import { UpstreamSpec as AwsUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { UpstreamSpec as AzureUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { UpstreamSpec as ConsulUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/consul/consul_pb';
import { UpstreamSpec as KubeUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/kubernetes/kubernetes_pb';
import {
  UpstreamSpec as Ec2UpstreamSpec,
  TagFilter
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/ec2/aws_ec2_pb';
import { UpstreamSpec as PipeUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/pipe/pipe_pb';
import { ServiceSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/service_spec_pb';
import { ServiceSpec as GrpcServiceSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/grpc/grpc_pb';
import { ServiceSpec as RestServiceSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/rest/rest_pb';

import {
  Host,
  UpstreamSpec as StaticUpstreamSpec
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/static/static_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { UpstreamApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb_service';
import { host } from 'store';
import { guardByLicense } from 'store/config/actions';
import { UPSTREAM_SPEC_TYPES } from 'utils/upstreamHelpers';
import {
  CreateUpstreamRequest,
  CreateUpstreamResponse,
  DeleteUpstreamRequest,
  DeleteUpstreamResponse,
  GetUpstreamRequest,
  GetUpstreamResponse,
  ListUpstreamsRequest,
  ListUpstreamsResponse,
  UpdateUpstreamRequest,
  UpdateUpstreamResponse,
  UpstreamInput,
  UpstreamDetails
} from '../../proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { Upstream } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/upstream_pb';
import { UpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import {
  UpstreamSslConfig,
  SSLFiles,
  SDSConfig,
  CallCredentials,
  SslParameters
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/ssl_pb';
import { CircuitBreakerConfig } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/circuit_breaker_pb';
import { LoadBalancerConfig } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/load_balancer_pb';
import { ConnectionConfig } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/connection_pb';
import { HealthCheck } from 'proto/github.com/solo-io/gloo/projects/gloo/api/external/envoy/api/v2/core/health_check_pb';
import { OutlierDetection } from 'proto/github.com/solo-io/gloo/projects/gloo/api/external/envoy/api/v2/cluster/outlier_detection_pb';
import { BoolValue } from 'google-protobuf/google/protobuf/wrappers_pb';
import { Metadata } from 'proto/github.com/solo-io/solo-kit/api/v1/metadata_pb';
import {
  TransformationTemplate,
  InjaTemplate
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/transformation/transformation_pb';
import {
  SubsetSpec,
  Selector
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/subset_spec_pb';

export const client = new UpstreamApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

// create
// update
// delete

function getUpstream(
  getUpstreamRequest: GetUpstreamRequest.AsObject
): Promise<UpstreamDetails> {
  return new Promise((resolve, reject) => {
    let req = new GetUpstreamRequest();
    let ref = new ResourceRef();
    ref.setName(getUpstreamRequest.ref!.name);
    ref.setNamespace(getUpstreamRequest.ref!.namespace);
    req.setRef(ref);

    client.getUpstream(req, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.getUpstreamDetails());
      }
    });
  });
}

function listUpstreams(): Promise<ListUpstreamsResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new ListUpstreamsRequest();

    client.listUpstreams(request, (error, data) => {
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

function setUpstreamValues(
  upstream: Upstream.AsObject,
  upstreamToUpdate = new Upstream()
): Upstream {
  let { upstreamSpec, status, metadata, discoveryMetadata } = upstream;

  if (metadata !== undefined) {
    let { name, namespace } = metadata;
    let newMetadata = new Metadata();
    newMetadata.setName(name);
    newMetadata.setNamespace(namespace);
    upstreamToUpdate.setMetadata(newMetadata);
  }

  if (upstreamSpec !== undefined) {
    let newUpstreamSpec = new UpstreamSpec();
    let {
      sslConfig,
      circuitBreakers,
      loadBalancerConfig,
      connectionConfig,
      healthChecksList,
      outlierDetection,
      useHttp2,
      kube,
      pb_static,
      pipe,
      aws,
      azure,
      consul,
      awsEc2
    } = upstreamSpec!;
    let newSslConfig = new UpstreamSslConfig();
    let newCircuitBreakers = new CircuitBreakerConfig();
    let newLoadBalancerConfig = new LoadBalancerConfig();
    let newConnectionConfig = new ConnectionConfig();
    let newOutlierDetection = new OutlierDetection();

    if (sslConfig !== undefined) {
      let {
        secretRef,
        sslFiles,
        sds,
        sni,
        verifySubjectAltNameList,
        parameters
      } = sslConfig!;

      if (secretRef !== undefined) {
        let sslSecretRef = new ResourceRef();
        sslSecretRef.setName(secretRef.name);
        sslSecretRef.setNamespace(secretRef.namespace);
      }
      // sslfiles
      if (sslFiles !== undefined) {
        let newSslFiles = new SSLFiles();
        let { rootCa, tlsCert, tlsKey } = sslFiles!;
        if (rootCa !== undefined) {
          newSslFiles.setRootCa(rootCa);
        }

        if (tlsCert !== undefined) {
          newSslFiles.setTlsCert(tlsCert);
        }

        if (tlsKey !== undefined) {
          newSslFiles.setTlsKey(tlsKey);
        }

        newSslConfig.setSslFiles(newSslFiles);
      }
      // sds
      if (sds !== undefined) {
        let newSdsConfig = new SDSConfig();
        let {
          targetUri,
          callCredentials,
          certificatesSecretName,
          validationContextName
        } = sds!;
        if (targetUri !== undefined) {
          newSdsConfig.setTargetUri(targetUri);
        }

        if (callCredentials !== undefined) {
          let newCallCreds = new CallCredentials();
          let { fileCredentialSource } = callCredentials!;
          if (fileCredentialSource !== undefined) {
            let newFileCredsSource = new CallCredentials.FileCredentialSource();
            let { tokenFileName, header } = fileCredentialSource!;
            newFileCredsSource.setHeader(header);
            newFileCredsSource.setTokenFileName(tokenFileName);

            newCallCreds.setFileCredentialSource(newFileCredsSource);
          }
          newSdsConfig.setCallCredentials(newCallCreds);
        }
        if (certificatesSecretName !== undefined) {
          newSdsConfig.setCertificatesSecretName(certificatesSecretName);
        }

        if (validationContextName !== undefined) {
          newSdsConfig.setValidationContextName(validationContextName);
        }
        newSslConfig.setSds(newSdsConfig);
      }

      //sni
      if (sni !== undefined) {
        newSslConfig.setSni(sni);
      }
      // verifysubkectaltnamelist
      if (verifySubjectAltNameList !== undefined) {
        newSslConfig.setVerifySubjectAltNameList(verifySubjectAltNameList);
      }
      // parameters
      if (parameters !== undefined) {
        let newSslParams = new SslParameters();
        let {
          maximumProtocolVersion,
          minimumProtocolVersion,
          cipherSuitesList,
          ecdhCurvesList
        } = parameters!;
        if (cipherSuitesList !== undefined) {
          newSslParams.setCipherSuitesList(cipherSuitesList);
        }
        if (maximumProtocolVersion !== undefined) {
          newSslParams.setMaximumProtocolVersion(maximumProtocolVersion);
        }

        if (minimumProtocolVersion !== undefined) {
          newSslParams.setMinimumProtocolVersion(minimumProtocolVersion);
        }

        if (ecdhCurvesList !== undefined) {
          newSslParams.setEcdhCurvesList(ecdhCurvesList);
        }
        newSslConfig.setParameters(newSslParams);
      }
      newUpstreamSpec.setSslConfig(newSslConfig);
    }
    if (circuitBreakers !== undefined) {
      // TODO
      newUpstreamSpec.setCircuitBreakers(newCircuitBreakers);
    }
    if (loadBalancerConfig !== undefined) {
      // TODO
      newUpstreamSpec.setLoadBalancerConfig(newLoadBalancerConfig);
    }

    if (connectionConfig !== undefined) {
      // TODO
      newUpstreamSpec.setConnectionConfig(newConnectionConfig);
    }

    // TODO
    if (healthChecksList !== undefined) {
      let newHealthChecksList = healthChecksList.map(healthCheck => {
        let newHealthCheck = new HealthCheck();
        return newHealthCheck;
      });
      newUpstreamSpec.setHealthChecksList(newHealthChecksList);
    }

    // TODO
    if (outlierDetection !== undefined) {
      newUpstreamSpec.setOutlierDetection(newOutlierDetection);
    }

    if (useHttp2 !== undefined) {
      newUpstreamSpec.setUseHttp2(useHttp2);
    }

    if (kube !== undefined) {
      let kubeSpec = new KubeUpstreamSpec();
      let {
        serviceName,
        serviceNamespace,
        servicePort,
        selectorMap,
        serviceSpec,
        subsetSpec
      } = kube!;
      if (serviceName !== undefined) {
        kubeSpec.setServiceName(serviceName);
      }
      if (serviceNamespace !== undefined) {
        kubeSpec.setServiceNamespace(serviceNamespace);
      }
      if (servicePort !== undefined) {
        kubeSpec.setServicePort(servicePort);
      }
      if (selectorMap !== undefined) {
        selectorMap.forEach(([key, val]) => {
          kubeSpec.getSelectorMap().set(key, val);
        });
      }
      if (serviceSpec !== undefined) {
        kubeSpec.setServiceSpec(setServiceSpecValues(serviceSpec));
      }
      if (subsetSpec !== undefined) {
        let newSubsetSpec = new SubsetSpec();
        let selList = subsetSpec.selectorsList.map(selector => {
          let newSelector = new Selector();
          newSelector.setKeysList(selector.keysList);
          return newSelector;
        });
        newSubsetSpec.setSelectorsList(selList);
        kubeSpec.setSubsetSpec(newSubsetSpec);
      }
      newUpstreamSpec.setKube(kubeSpec);
    }

    if (pb_static !== undefined) {
      let newStatic = new StaticUpstreamSpec();
      let { hostsList, useTls, serviceSpec } = pb_static!;
      let newHostsList = hostsList.map(host => {
        let newHost = new Host();
        newHost.setAddr(host.addr);
        newHost.setPort(host.port);
        return newHost;
      });
      newStatic.setHostsList(newHostsList);

      if (useTls !== undefined) {
        newStatic.setUseTls(useTls);
      }

      if (serviceSpec !== undefined) {
        newStatic.setServiceSpec(setServiceSpecValues(serviceSpec));
      }

      newUpstreamSpec.setStatic(newStatic);
    }

    if (pipe !== undefined) {
      let { path, serviceSpec } = pipe!;
      let newPipeUpstreamSpec = new PipeUpstreamSpec();
      if (path !== undefined) {
        newPipeUpstreamSpec.setPath(path);
      }
      if (serviceSpec !== undefined) {
        let newPipeServiceSpec = setServiceSpecValues(serviceSpec);
        newPipeUpstreamSpec.setServiceSpec(newPipeServiceSpec);
      }
      newUpstreamSpec.setPipe(newPipeUpstreamSpec);
    }

    if (aws !== undefined) {
      let { region, secretRef, lambdaFunctionsList } = aws!;
      let newAwsSpec = new AwsUpstreamSpec();
      if (region !== undefined) {
        newAwsSpec.setRegion(region);
      }
      if (secretRef !== undefined) {
        let awsSecretRef = new ResourceRef();
        awsSecretRef.setName(secretRef!.name);
        awsSecretRef.setNamespace(secretRef!.namespace);
        newAwsSpec.setSecretRef(awsSecretRef);
      }
      newUpstreamSpec.setAws(newAwsSpec);
    }

    if (azure !== undefined) {
      let azureSpec = new AzureUpstreamSpec();
      let { functionAppName, secretRef } = azure!;
      if (secretRef !== undefined) {
        let azureRef = new ResourceRef();
        azureRef.setName(secretRef.name);
        azureRef.setNamespace(secretRef.namespace);
        azureSpec.setSecretRef(azureRef);
      }
      if (functionAppName !== undefined) {
        azureSpec.setFunctionAppName(functionAppName);
      }
      newUpstreamSpec.setAzure(azureSpec);
    }

    if (consul !== undefined) {
      let consulSpec = new ConsulUpstreamSpec();
      let {
        serviceName,
        serviceTagsList,
        serviceSpec,
        connectEnabled,
        dataCentersList
      } = consul!;
      if (serviceName !== undefined) {
        consulSpec.setServiceName(serviceName);
      }
      if (serviceTagsList !== undefined) {
        consulSpec.setServiceTagsList(serviceTagsList);
      }

      if (serviceSpec !== undefined) {
        consulSpec.setServiceSpec(setServiceSpecValues(serviceSpec));
      }
      if (connectEnabled !== undefined) {
        consulSpec.setConnectEnabled(connectEnabled);
      }
      if (dataCentersList !== undefined) {
        consulSpec.setDataCentersList(dataCentersList);
      }

      newUpstreamSpec.setConsul(consulSpec);
    }

    if (awsEc2 !== undefined) {
      let ec2Spec = new Ec2UpstreamSpec();
      let {
        region,
        secretRef,
        roleArn,
        roleArnsList, //deprecated
        filtersList,
        publicIp,
        port
      } = awsEc2!;
      if (region !== undefined) {
        ec2Spec.setRegion(region);
      }

      if (secretRef !== undefined) {
        let ec2Secret = new ResourceRef();
        ec2Secret.setName(secretRef.name);
        ec2Secret.setNamespace(secretRef.namespace);
      }

      if (roleArn !== undefined) {
        ec2Spec.setRoleArn(roleArn);
      }

      if (filtersList !== undefined) {
        let filters = filtersList.map(filt => {
          let { key, kvPair } = filt;
          let newTagFilter = new TagFilter();
          if (key !== undefined) {
            newTagFilter.setKey(key);
          } else if (kvPair !== undefined) {
            let kv = new TagFilter.KvPair();
            kv.setKey(kvPair.key);
            kv.setValue(kvPair.value);
            newTagFilter.setKvPair(kv);
          }
          return newTagFilter;
        });
        ec2Spec.setFiltersList(filters);
      }
      if (publicIp !== undefined) {
        ec2Spec.setPublicIp(publicIp);
      }

      if (port !== undefined) {
        ec2Spec.setPort(port);
      }
      newUpstreamSpec.setAwsEc2(ec2Spec);
    }
    upstreamToUpdate.setUpstreamSpec(newUpstreamSpec);
  }
  return upstreamToUpdate;
}

function setServiceSpecValues(serviceSpec: ServiceSpec.AsObject): ServiceSpec {
  let serviceSpecToSet = new ServiceSpec();
  let { grpc, rest } = serviceSpec!;
  if (grpc !== undefined) {
    let grpcSpec = new GrpcServiceSpec();
    let { descriptors, grpcServicesList } = grpc!;
    if (descriptors !== undefined) {
      grpcSpec.setDescriptors(descriptors);
    }
    if (grpcServicesList !== undefined) {
      let grpcServices = grpcServicesList.map(grpc => {
        let grpcService = new GrpcServiceSpec.GrpcService();
        grpcService.setPackageName(grpc.packageName);
        grpcService.setServiceName(grpc.serviceName);
        grpcService.setFunctionNamesList(grpc.functionNamesList);
        return grpcService;
      });
      grpcSpec.setGrpcServicesList(grpcServices);
    }
    serviceSpecToSet.setGrpc(grpcSpec);
  }
  if (rest !== undefined) {
    let restSpec = new RestServiceSpec();
    let { transformationsMap, swaggerInfo } = rest!;
    if (swaggerInfo !== undefined) {
      let swagger = new RestServiceSpec.SwaggerInfo();
      let { url, inline } = swaggerInfo!;
      if (url !== undefined) {
        swagger.setUrl(url);
      }
      if (inline !== undefined) {
        swagger.setInline(inline);
      }
      restSpec.setSwaggerInfo(swagger);
    }
    if (transformationsMap !== undefined) {
      let transformMap = restSpec.getTransformationsMap();

      transformationsMap.forEach(([transformName, template]) => {
        let temp = new TransformationTemplate();
        temp.setAdvancedTemplates(template.advancedTemplates);
        let templateBody = new InjaTemplate();
        templateBody.setText(template.body!.text);
        temp.setBody(templateBody);

        transformMap.set(transformName, temp);
      });
    }
    serviceSpecToSet.setRest(restSpec);
  }
  return serviceSpecToSet;
}

function createUpstream(
  createUpstreamRequest: CreateUpstreamRequest.AsObject
): Promise<CreateUpstreamResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new CreateUpstreamRequest();
    const { upstreamInput } = createUpstreamRequest;
    if (upstreamInput !== undefined) {
      let inputUpstream = setUpstreamValues(upstreamInput);

      request.setUpstreamInput(inputUpstream);
    }

    client.createUpstream(request, (error, data) => {
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

function updateUpstream(
  updateUpstreamRequest: UpdateUpstreamRequest.AsObject
): Promise<UpdateUpstreamResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let request = new UpdateUpstreamRequest();
    let { upstreamInput } = updateUpstreamRequest!;
    if (upstreamInput !== undefined && upstreamInput.metadata !== undefined) {
      let { name, namespace } = upstreamInput.metadata;
      let upstreamToUpdate = await getUpstream({
        ref: {
          name,
          namespace
        }
      });
      let updatedUpstream = setUpstreamValues(
        upstreamInput,
        upstreamToUpdate.getUpstream()
      );
      request.setUpstreamInput(updatedUpstream);
    }

    client.updateUpstream(request, (error, data) => {
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

function deleteUpstream(
  deleteUpstreamRequest: DeleteUpstreamRequest.AsObject
): Promise<DeleteUpstreamResponse> {
  return new Promise((resolve, reject) => {
    let request = new DeleteUpstreamRequest();
    let ref = new ResourceRef();
    ref.setName(deleteUpstreamRequest.ref!.name);
    ref.setNamespace(deleteUpstreamRequest.ref!.namespace);
    request.setRef(ref);
    guardByLicense();
    client.deleteUpstream(request, (error, data) => {
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

export const upstreams = {
  getUpstream,
  listUpstreams,
  createUpstream,
  updateUpstream,
  deleteUpstream
};
