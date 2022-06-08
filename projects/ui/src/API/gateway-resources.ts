import { GatewayResourceApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb_service';
import {
  host,
  getObjectRefClassFromRefObj,
  toPaginationClass,
  getClusterRefClassFromClusterRefObj,
} from './helpers';
import { grpc } from '@improbable-eng/grpc-web';
import {
  VirtualService,
  ListVirtualServicesRequest,
  ListGatewaysRequest,
  Gateway,
  GetVirtualServiceYamlRequest,
  GetGatewayYamlRequest,
  GetRouteTableYamlRequest,
  ListRouteTablesRequest,
  ListRouteTablesResponse,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';

import {
  Pagination,
  StatusFilter,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb';

import {
  ObjectRef,
  ClusterObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

const gatewayResourceApiClient = new GatewayResourceApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true,
});

export const gatewayResourceApi = {
  listVirtualServices,
  listRouteTables,
  listGateways,
  getVirtualServiceYAML,
  getGatewayYAML,
  getRouteTableYAML,
};

function listVirtualServices(
  listVSRequest?: ObjectRef.AsObject
): Promise<VirtualService.AsObject[]> {
  let request = new ListVirtualServicesRequest();
  if (listVSRequest) {
    request.setGlooInstanceRef(getObjectRefClassFromRefObj(listVSRequest));
  }

  return new Promise((resolve, reject) => {
    gatewayResourceApiClient.listVirtualServices(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().virtualServicesList);
      }
    });
  });
}

function listRouteTables(
  listRouteTableRequest?: ObjectRef.AsObject,
  pagination?: Pagination.AsObject,
  queryString?: string,
  statusFilter?: number
): Promise<ListRouteTablesResponse.AsObject> {
  let request = new ListRouteTablesRequest();
  if (listRouteTableRequest) {
    request.setGlooInstanceRef(
      getObjectRefClassFromRefObj(listRouteTableRequest)
    );
  }
  if (pagination) {
    request.setPagination(toPaginationClass(pagination));
  }
  if (queryString) {
    request.setQueryString(queryString);
  }
  if (statusFilter !== undefined) {
    const sf = new StatusFilter();
    sf.setState(statusFilter);
    request.setStatusFilter(sf);
  }

  return new Promise((resolve, reject) => {
    gatewayResourceApiClient.listRouteTables(request, (error, data) => {
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

function listGateways(
  listGatewayRequest?: ObjectRef.AsObject
): Promise<Gateway.AsObject[]> {
  let request = new ListGatewaysRequest();
  if (listGatewayRequest) {
    request.setGlooInstanceRef(getObjectRefClassFromRefObj(listGatewayRequest));
  }

  return new Promise((resolve, reject) => {
    gatewayResourceApiClient.listGateways(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().gatewaysList);
      }
    });
  });
}

function getVirtualServiceYAML(
  vsObjectRef: ClusterObjectRef.AsObject
): Promise<string> {
  let request = new GetVirtualServiceYamlRequest();
  request.setVirtualServiceRef(
    getClusterRefClassFromClusterRefObj(vsObjectRef)
  );

  return new Promise((resolve, reject) => {
    gatewayResourceApiClient.getVirtualServiceYaml(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().yamlData?.yaml ?? 'None');
      }
    });
  });
}

function getGatewayYAML(
  gatewayObjectRef: ClusterObjectRef.AsObject
): Promise<string> {
  let request = new GetGatewayYamlRequest();
  request.setGatewayRef(getClusterRefClassFromClusterRefObj(gatewayObjectRef));

  return new Promise((resolve, reject) => {
    gatewayResourceApiClient.getGatewayYaml(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().yamlData?.yaml ?? 'None');
      }
    });
  });
}

function getRouteTableYAML(
  rtObjectRef: ClusterObjectRef.AsObject
): Promise<string> {
  let request = new GetRouteTableYamlRequest();
  request.setRouteTableRef(getClusterRefClassFromClusterRefObj(rtObjectRef));

  return new Promise((resolve, reject) => {
    gatewayResourceApiClient.getRouteTableYaml(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().yamlData?.yaml ?? 'None');
      }
    });
  });
}
