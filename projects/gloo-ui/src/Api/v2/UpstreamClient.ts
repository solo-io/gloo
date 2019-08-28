import {
  ListUpstreamsRequest,
  ListUpstreamsResponse,
  GetUpstreamRequest,
  CreateUpstreamRequest,
  UpdateUpstreamRequest,
  DeleteUpstreamRequest,
  GetUpstreamResponse,
  CreateUpstreamResponse,
  UpdateUpstreamResponse,
  DeleteUpstreamResponse,
  UpstreamInput
} from '../../proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { UpstreamApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb_service';
import { UpstreamSpec as AwsUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/aws/aws_pb';
import { host } from 'store';
import { grpc } from '@improbable-eng/grpc-web';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { UPSTREAM_SPEC_TYPES } from 'utils/upstreamHelpers';
import { AwsValuesType } from 'Components/Features/Upstream/Creation/AwsUpstreamForm';
import { KubeValuesType } from 'Components/Features/Upstream/Creation/KubeUpstreamForm';
import { StaticValuesType } from 'Components/Features/Upstream/Creation/StaticUpstreamForm';
import { AzureValuesType } from 'Components/Features/Upstream/Creation/AzureUpstreamForm';
import { ConsulVauesType } from 'Components/Features/Upstream/Creation/ConsulUpstreamForm';
import { UpstreamSpec as AzureUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { UpstreamSpec as ConsulUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/consul/consul_pb';
import { UpstreamSpec as KubeUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/kubernetes/kubernetes_pb';
import { ServiceSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/service_spec_pb';
import {
  UpstreamSpec as StaticUpstreamSpec,
  Host
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/static/static_pb';

export const client = new UpstreamApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

export interface UpstreamSpecificValues
  extends AwsValuesType,
    KubeValuesType,
    StaticValuesType,
    AzureValuesType,
    ConsulVauesType {}

function getUpstreamsList(
  listUpstreamsRequest: ListUpstreamsRequest.AsObject
): Promise<ListUpstreamsResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let req = new ListUpstreamsRequest();
    req.setNamespacesList(listUpstreamsRequest.namespacesList);

    client.listUpstreams(req, (error, data) => {
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

function getUpstream(
  getUpstreamRequest: GetUpstreamRequest.AsObject
): Promise<GetUpstreamResponse.AsObject> {
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
        resolve(data!.toObject());
      }
    });
  });
}

function createUpstream(params: {
  name: string;
  namespace: string;
  type: string;
  values: UpstreamSpecificValues;
}): Promise<CreateUpstreamResponse> {
  const { values, name, namespace, type } = params;
  return new Promise((resolve, reject) => {
    let req = new CreateUpstreamRequest();
    let upstreamInput = getUpstreamInput({ name, namespace, type, values });
    req.setInput(upstreamInput);

    client.createUpstream(req, (error, data) => {
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

function getUpstreamInput(params: {
  name: string;
  namespace: string;
  type: string;
  values: UpstreamSpecificValues;
}): UpstreamInput {
  const { name, namespace, type, values } = params;
  let newUpstream = new UpstreamInput();
  // set the resource ref
  let ref = new ResourceRef();
  ref.setName(name);
  ref.setNamespace(namespace);
  newUpstream.setRef(ref);
  // upstream specific values
  switch (type) {
    case UPSTREAM_SPEC_TYPES.AWS:
      const awsSpec = new AwsUpstreamSpec();
      awsSpec.setRegion(values.awsRegion);
      const awsSecretRef = new ResourceRef();
      awsSecretRef.setName(values.awsSecretRef.name);
      awsSecretRef.setNamespace(values.awsSecretRef.namespace);
      awsSpec.setSecretRef(awsSecretRef);
      newUpstream.setAws(awsSpec);
      break;
    case UPSTREAM_SPEC_TYPES.AZURE:
      const azureSpec = new AzureUpstreamSpec();
      const azureSecretRef = new ResourceRef();
      azureSecretRef.setName(values.azureSecretRef.name);
      azureSecretRef.setNamespace(values.azureSecretRef.namespace);
      azureSpec.setSecretRef(azureSecretRef);
      azureSpec.setFunctionAppName(values.azureFunctionAppName);
      newUpstream.setAzure(azureSpec);
      break;
    case UPSTREAM_SPEC_TYPES.KUBE:
      const kubeSpec = new KubeUpstreamSpec();
      kubeSpec.setServiceName(values.kubeServiceName);
      kubeSpec.setServiceNamespace(values.kubeServiceNamespace);
      kubeSpec.setServicePort(values.kubeServicePort);
      newUpstream.setKube(kubeSpec);
      break;
    case UPSTREAM_SPEC_TYPES.STATIC:
      const staticSpec = new StaticUpstreamSpec();
      staticSpec.setUseTls(values.staticUseTls);
      newUpstream.setStatic(staticSpec);
      break;
    case UPSTREAM_SPEC_TYPES.CONSUL:
      const consulSpec = new ConsulUpstreamSpec();
      consulSpec.setServiceName(values.consulServiceName);
      consulSpec.setServiceTagsList(values.consulServiceTagsList);
      consulSpec.setConnectEnabled(values.consulConnectEnabled);
      consulSpec.setDataCentersList(values.consulDataCentersList);
      const consulServiceSpec = new ServiceSpec();
      consulSpec.setServiceSpec(consulServiceSpec);
      newUpstream.setConsul(consulSpec);
      break;
    default:
      throw new Error('not supported');
  }
  return newUpstream;
}

function updateUpstream(params: {
  name: string;
  namespace: string;
  type: string;
  values: UpstreamSpecificValues;
}): Promise<UpdateUpstreamResponse> {
  const { name, namespace, type, values } = params;
  let updateReq = new UpdateUpstreamRequest();
  updateReq.setInput(getUpstreamInput({ name, namespace, type, values }));
  return new Promise((resolve, reject) => {
    client.updateUpstream(updateReq, (error, data) => {
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

function deleteUpstream(
  deleteUpstreamRequest: DeleteUpstreamRequest.AsObject
): Promise<DeleteUpstreamResponse> {
  return new Promise((resolve, reject) => {
    let request = new DeleteUpstreamRequest();
    let ref = new ResourceRef();
    ref.setName(deleteUpstreamRequest.ref!.name);
    ref.setNamespace(deleteUpstreamRequest.ref!.namespace);
    request.setRef(ref);
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

export function getCreateUpstream(
  createUpstreamRequest: CreateUpstreamRequest.AsObject
): Promise<CreateUpstreamResponse.AsObject> {
  return new Promise((resolve, reject) => {
    const { input } = createUpstreamRequest;
    let ref = new ResourceRef();
    let request = new CreateUpstreamRequest();
    let usInput = new UpstreamInput();

    request.setInput();
    ref.setName(input!.ref!.name);
    ref.setNamespace(input!.ref!.namespace);

    usInput.setRef(ref);
    let awsSpec = new AwsUpstreamSpec();
    let azureSpec = new AzureUpstreamSpec();
    let staticSpec = new StaticUpstreamSpec();
    let kubeSpec = new KubeUpstreamSpec();
    let consulSpec = new ConsulUpstreamSpec();

    if (input!.aws) {
      const { region, secretRef } = input!.aws;
      let awsSecretRef = new ResourceRef();
      awsSecretRef.setName(secretRef!.name);
      awsSecretRef.setNamespace(secretRef!.namespace);
      awsSpec.setRegion(region);
      awsSpec.setSecretRef(awsSecretRef);
      usInput.setAws(awsSpec);
    } else if (input!.pb_static) {
      const { useTls, hostsList /*serviceSpec*/ } = input!.pb_static!;
      staticSpec.setUseTls(useTls);
      let hosts = hostsList.map(host => {
        let hostAdded = new Host();
        hostAdded.setAddr(host.addr);
        hostAdded.setPort(host.port);
        return hostAdded;
      });
      staticSpec.setHostsList(hosts);
      usInput.setStatic(staticSpec);
    } else if (input!.kube) {
      const { serviceName, serviceNamespace, servicePort } = input!.kube!;
      kubeSpec.setServiceName(serviceName);
      kubeSpec.setServiceNamespace(serviceNamespace);
      kubeSpec.setServicePort(servicePort);
      usInput.setKube(kubeSpec);
    } else if (input!.azure) {
      const { functionAppName, secretRef } = input!.azure!;
      const azureSecretRef = new ResourceRef();
      azureSecretRef.setName(secretRef!.name);
      azureSecretRef.setNamespace(secretRef!.namespace);
      azureSpec.setSecretRef(azureSecretRef);
      azureSpec.setFunctionAppName(functionAppName);
      usInput.setAzure(azureSpec);
    } else if (input!.consul) {
      const {
        connectEnabled,
        dataCentersList,
        serviceName,
        //serviceSpec,
        serviceTagsList
      } = input!.consul!;
      consulSpec.setServiceName(serviceName);
      consulSpec.setServiceTagsList(serviceTagsList);
      consulSpec.setConnectEnabled(connectEnabled);
      consulSpec.setDataCentersList(dataCentersList);
      const consulServiceSpec = new ServiceSpec();
      consulSpec.setServiceSpec(consulServiceSpec);
      usInput.setConsul(consulSpec);
    }
    request.setInput(usInput);
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

export const upstreams = {
  getUpstream,
  getUpstreamsList,
  getCreateUpstream,
  createUpstream,
  updateUpstream,
  deleteUpstream
};
