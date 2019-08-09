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
import { host } from '../grpc-web-hooks';
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
import { UpstreamSpec as StaticUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/static/static_pb';

const client = new UpstreamApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getUpstreamsList(params: {
  namespaces: string[];
}): Promise<ListUpstreamsResponse> {
  return new Promise((resolve, reject) => {
    let req = new ListUpstreamsRequest();
    req.setNamespacesList(params.namespaces);
    client.listUpstreams(req, (error, data) => {
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

function getUpstream(params: {
  name: string;
  namespace: string;
}): Promise<GetUpstreamResponse> {
  return new Promise((resolve, reject) => {
    let req = new GetUpstreamRequest();
    let ref = new ResourceRef();
    ref.setName(params.name);
    ref.setNamespace(params.namespace);
    req.setRef(ref);
    client.getUpstream(req, (error, data) => {
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

export interface UpstreamSpecificValues
  extends AwsValuesType,
    KubeValuesType,
    StaticValuesType,
    AzureValuesType,
    ConsulVauesType {}

function createUpstream(params: {
  name: string;
  namespace: string;
  type: string;
  values: UpstreamSpecificValues;
}): Promise<CreateUpstreamResponse> {
  const { values, name, namespace, type } = params;
  return new Promise((resolve, reject) => {
    let req = new CreateUpstreamRequest();
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
      default:
        break;
    }

    req.setInput(newUpstream);
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

// TODO
function updateUpstream(params: {}) {}

function deleteUpstream(params: {
  name: string;
  namespace: string;
}): Promise<DeleteUpstreamResponse> {
  return new Promise((resolve, reject) => {
    let request = new DeleteUpstreamRequest();
    let ref = new ResourceRef();
    ref.setName(params.name);
    ref.setNamespace(params.namespace);
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

export const upstreams = {
  getUpstream,
  getUpstreamsList,
  createUpstream,
  updateUpstream,
  deleteUpstream
};
