import { ProxyApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/proxy_pb_service';
import { host } from '../grpc-web-hooks';
import { grpc } from '@improbable-eng/grpc-web';
import {
  ListProxiesResponse,
  ListProxiesRequest
} from '../../proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/proxy_pb';

export const client = new ProxyApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getProxiesList(params: {
  namespaces: string[];
}): Promise<ListProxiesResponse> {
  return new Promise((resolve, reject) => {
    let req = new ListProxiesRequest();
    req.setNamespacesList(params.namespaces);
    client.listProxies(req, (error, data) => {
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

export const proxy = {
  getProxiesList
};
