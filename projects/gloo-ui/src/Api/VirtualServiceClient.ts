// virtualservice client goes here
import { grpc } from '@improbable-eng/grpc-web';
import { VirtualServiceApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb_service';
import { host } from './grpc-web-hooks';
import {
  ListVirtualServicesResponse,
  ListVirtualServicesRequest,
  GetVirtualServiceResponse,
  GetVirtualServiceRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';

export const client = new VirtualServiceApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({
    withCredentials: false
  }),
  debug: true
});

function getVirtualServicesList(params: {
  namespaces: string[];
}): Promise<ListVirtualServicesResponse> {
  return new Promise((resolve, reject) => {
    let request = new ListVirtualServicesRequest();
    request.setNamespacesList(params.namespaces);
    client.listVirtualServices(request, (error, data) => {
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

function getVirtualService(params: {
  name: string;
  namespace: string;
}): Promise<GetVirtualServiceResponse> {
  return new Promise((resolve, reject) => {
    let request = new GetVirtualServiceRequest();
    let ref = new ResourceRef();
    ref.setName(params.name);
    ref.setNamespace(params.namespace);
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

export const virtualServices = {
  getVirtualService,
  getVirtualServicesList
};
