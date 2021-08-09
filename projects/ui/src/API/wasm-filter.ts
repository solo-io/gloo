import { WasmFilterApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/wasm_pb_service';
import { getObjectRefClassFromRefObj, host } from './helpers';
import { grpc } from '@improbable-eng/grpc-web';
import {
  WasmFilter,
  ListWasmFiltersRequest,
  DescribeWasmFilterRequest,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/wasm_pb';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

/*****
 * WARNING!! This API is still supported, but does NOT currently power anything on the front-end.
 *   If you are looking for information on Routes that the front-end uses, you are most likely looking
 *   for the VirtualServiceRoutesApi in the similarly naamed ts file.
 */

const wasmFilterApiClient = new WasmFilterApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true,
});

export const wasmFilterApi = {
  getWasmFilter,
  listWasmFilters,
};

function getWasmFilter(
  getWasmFilterRequest?: DescribeWasmFilterRequest.AsObject
): Promise<WasmFilter.AsObject> {
  let request = new DescribeWasmFilterRequest();
  if (getWasmFilterRequest) {
    request.setName(getWasmFilterRequest.name);
    if (getWasmFilterRequest.gatewayRef) {
      request.setGatewayRef(
        getObjectRefClassFromRefObj(getWasmFilterRequest.gatewayRef)
      );
    }
    request.setRootId(getWasmFilterRequest.rootId);
  }

  return new Promise((resolve, reject) => {
    wasmFilterApiClient.describeWasmFilter(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        const filter = data!.toObject().wasmFilter;
        if (filter !== undefined) resolve(filter);
        else {
          reject({ message: 'Wasm filter not found' });
        }
      }
    });
  });
}

function listWasmFilters(): Promise<WasmFilter.AsObject[]> {
  let request = new ListWasmFiltersRequest();

  return new Promise((resolve, reject) => {
    wasmFilterApiClient.listWasmFilters(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().wasmFiltersList);
      }
    });
  });
}
