import { VirtualServiceRoutesApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/rt_selector_pb_service';
import { host, getClusterRefClassFromClusterRefObj } from './helpers';
import { grpc } from '@improbable-eng/grpc-web';
import {
  GetVirtualServiceRoutesRequest,
  SubRouteTableRow
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/rt_selector_pb';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

/*****
 * WARNING!! This API is still supported, but does NOT currently power anything on the front-end.
 *   If you are looking for information on Routes that the front-end uses, you are most likely looking
 *   for the VirtualServiceRoutesApi in the similarly naamed ts file.
 */

const virtualServiceRoutesApiClient = new VirtualServiceRoutesApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

export const routeTablesSelectorApi = {
  getSubroutesForVirtualService
};

function getSubroutesForVirtualService(
  listRoutesRequest?: ClusterObjectRef.AsObject
): Promise<SubRouteTableRow.AsObject[]> {
  let request = new GetVirtualServiceRoutesRequest();
  if (listRoutesRequest) {
    request.setVirtualServiceRef(
      getClusterRefClassFromClusterRefObj(listRoutesRequest)
    );
  }

  return new Promise((resolve, reject) => {
    virtualServiceRoutesApiClient.getVirtualServiceRoutes(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().vsRoutesList);
        }
      }
    );
  });
}
