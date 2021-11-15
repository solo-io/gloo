import { GlooInstanceApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb_service';
import { getObjectRefClassFromRefObj, host } from './helpers';
import { grpc } from '@improbable-eng/grpc-web';
import {
  GlooInstance,
  ListGlooInstancesRequest,
  ClusterDetails,
  ListClusterDetailsRequest,
  GetConfigDumpsRequest,
  ConfigDump,
  GetUpstreamHostsRequest,
  HostList,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { ObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

const glooInstanceApiClient = new GlooInstanceApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true,
});

export const glooInstanceApi = {
  listGlooInstances,
  listClusterDetails,
  getConfigDumps,
  getUpstreamHosts,
};

function listGlooInstances(): Promise<GlooInstance.AsObject[]> {
  let request = new ListGlooInstancesRequest();

  return new Promise((resolve, reject) => {
    glooInstanceApiClient.listGlooInstances(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().glooInstancesList);
      }
    });
  });
}

function listClusterDetails(): Promise<ClusterDetails.AsObject[]> {
  let request = new ListClusterDetailsRequest();

  return new Promise((resolve, reject) => {
    glooInstanceApiClient.listClusterDetails(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().clusterDetailsList);
      }
    });
  });
}

function getConfigDumps(
  instanceRef: ObjectRef.AsObject
): Promise<ConfigDump.AsObject[]> {
  let request = new GetConfigDumpsRequest();
  request.setGlooInstanceRef(getObjectRefClassFromRefObj(instanceRef));

  return new Promise((resolve, reject) => {
    glooInstanceApiClient.getConfigDumps(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().configDumpsList);
      }
    });
  });
}

function getUpstreamHosts(
  instanceRef: ObjectRef.AsObject
): Promise<Map<string, HostList.AsObject>> {
  let request = new GetUpstreamHostsRequest();
  request.setGlooInstanceRef(getObjectRefClassFromRefObj(instanceRef));

  return new Promise((resolve, reject) => {
    glooInstanceApiClient.getUpstreamHosts(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        // for some reason the protobuf map gets compiled into an array in typescript;
        // turn it back into a map here
        const upstreamHostArray = data!.toObject().upstreamHostsMap;
        const upstreamHostMap: Map<string, HostList.AsObject> = new Map<
          string,
          HostList.AsObject
        >();
        upstreamHostArray.forEach(([upstreamKey, hostList]) => {
          upstreamHostMap.set(upstreamKey, hostList);
        });
        resolve(upstreamHostMap);
      }
    });
  });
}
