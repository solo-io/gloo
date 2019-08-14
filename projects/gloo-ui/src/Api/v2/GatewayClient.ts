import { GatewayApiClient } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb_service';
import { host } from '../grpc-web-hooks';
import { grpc } from '@improbable-eng/grpc-web';
import { ListGatewaysResponse, ListGatewaysRequest } from '../../proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';

const client = new GatewayApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

function getGatewaysList(params: {
  namespaces: string[];
}): Promise<ListGatewaysResponse> {
  return new Promise((resolve, reject) => {
    let req = new ListGatewaysRequest();
    req.setNamespacesList(params.namespaces)
    client.listGateways(req, (error, data) => {
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

export const gateways = {
  getGatewaysList
};
