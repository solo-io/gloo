import { FederatedGatewayResourceApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb_service';
import { host, getObjectRefClassFromRefObj } from './helpers';
import { grpc } from '@improbable-eng/grpc-web';
import {
  FederatedVirtualService,
  FederatedRouteTable,
  FederatedGateway,
  ListFederatedVirtualServicesRequest,
  ListFederatedRouteTablesRequest,
  ListFederatedGatewaysRequest,
  GetFederatedVirtualServiceYamlRequest,
  GetFederatedRouteTableYamlRequest,
  GetFederatedGatewayYamlRequest,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb';
import { ObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

const federatedGatewayResourceApiClient = new FederatedGatewayResourceApiClient(
  host,
  {
    transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
    debug: true,
  }
);

export const federatedGatewayResourceApi = {
  listFederatedVirtualServices,
  listFederatedRouteTables,
  listFederatedGateways,
  getFederatedVirtualServiceYAML,
  getFederatedRouteTableYAML,
  getFederatedGatewayYAML,
};

function listFederatedVirtualServices(): Promise<
  FederatedVirtualService.AsObject[]
> {
  let request = new ListFederatedVirtualServicesRequest();

  return new Promise((resolve, reject) => {
    federatedGatewayResourceApiClient.listFederatedVirtualServices(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().federatedVirtualServicesList);
        }
      }
    );
  });
}

function listFederatedRouteTables(): Promise<FederatedRouteTable.AsObject[]> {
  let request = new ListFederatedRouteTablesRequest();

  return new Promise((resolve, reject) => {
    federatedGatewayResourceApiClient.listFederatedRouteTables(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().federatedRouteTablesList);
        }
      }
    );
  });
}

function listFederatedGateways(): Promise<FederatedGateway.AsObject[]> {
  let request = new ListFederatedGatewaysRequest();

  return new Promise((resolve, reject) => {
    federatedGatewayResourceApiClient.listFederatedGateways(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().federatedGatewaysList);
        }
      }
    );
  });
}

function getFederatedVirtualServiceYAML(
  fedVSObjectRef: ObjectRef.AsObject
): Promise<string> {
  let request = new GetFederatedVirtualServiceYamlRequest();
  request.setFederatedVirtualServiceRef(
    getObjectRefClassFromRefObj(fedVSObjectRef)
  );

  return new Promise((resolve, reject) => {
    federatedGatewayResourceApiClient.getFederatedVirtualServiceYaml(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().yamlData?.yaml ?? 'None');
        }
      }
    );
  });
}

function getFederatedRouteTableYAML(
  fedRTObjectRef: ObjectRef.AsObject
): Promise<string> {
  let request = new GetFederatedRouteTableYamlRequest();
  request.setFederatedRouteTableRef(
    getObjectRefClassFromRefObj(fedRTObjectRef)
  );

  return new Promise((resolve, reject) => {
    federatedGatewayResourceApiClient.getFederatedRouteTableYaml(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().yamlData?.yaml ?? 'None');
        }
      }
    );
  });
}

function getFederatedGatewayYAML(
  fedGatewayObjectRef: ObjectRef.AsObject
): Promise<string> {
  let request = new GetFederatedGatewayYamlRequest();
  request.setFederatedGatewayRef(
    getObjectRefClassFromRefObj(fedGatewayObjectRef)
  );

  return new Promise((resolve, reject) => {
    federatedGatewayResourceApiClient.getFederatedGatewayYaml(
      request,
      (error, data) => {
        if (error !== null) {
          console.error('Error:', error.message);
          console.error('Code:', error.code);
          console.error('Metadata:', error.metadata);
          reject(error);
        } else {
          resolve(data!.toObject().yamlData?.yaml ?? 'None');
        }
      }
    );
  });
}
