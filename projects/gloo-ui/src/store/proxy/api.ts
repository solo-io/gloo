import {
  ListProxiesRequest,
  ListProxiesResponse,
  ProxyDetails
} from 'proto/solo-projects/projects/grpcserver/api/v1/proxy_pb';
import { ProxyApiClient } from 'proto/solo-projects/projects/grpcserver/api/v1/proxy_pb_service';
import { host } from 'store';
import { grpc } from '@improbable-eng/grpc-web';

const client = new ProxyApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getListProxies(): Promise<ProxyDetails.AsObject[]> {
  return new Promise((resolve, reject) => {
    let request = new ListProxiesRequest();
    client.listProxies(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().proxyDetailsList);
      }
    });
  });
}

export const proxyAPI = {
  getListProxies
};
