import { GlooResourceApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb_service';
import {
  host,
  getObjectRefClassFromRefObj,
  getClusterRefClassFromClusterRefObj,
  toPaginationClass,
} from './helpers';
import { grpc } from '@improbable-eng/grpc-web';
import {
  ListUpstreamGroupsRequest,
  ListUpstreamsRequest,
  ListUpstreamsResponse,
  UpstreamGroup,
  Upstream,
  GetUpstreamGroupYamlRequest,
  GetUpstreamYamlRequest,
  GetProxyYamlRequest,
  ListProxiesRequest,
  Proxy,
  Settings,
  GetSettingsYamlRequest,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import {
  ObjectRef,
  ClusterObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import {
  Pagination,
  StatusFilter,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb';

const glooResourceApiClient = new GlooResourceApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true,
});

export const glooResourceApi = {
  listUpstreams,
  listUpstreamGroups,
  listProxies,
  listSettings,
  getUpstream,
  getUpstreamGroup,
  getUpstreamYAML,
  getUpstreamGroupYAML,
  getProxyYAML,
  getSettingYAML,
};

function listUpstreams(
  listUpstreamsRequest?: ObjectRef.AsObject,
  pagination?: Pagination.AsObject,
  queryString?: string,
  statusFilter?: number
): Promise<ListUpstreamsResponse.AsObject> {
  let request = new ListUpstreamsRequest();
  if (listUpstreamsRequest) {
    request.setGlooInstanceRef(
      getObjectRefClassFromRefObj(listUpstreamsRequest)
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
    glooResourceApiClient.listUpstreams(request, (error, data) => {
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

function listUpstreamGroups(
  listUpstreamGroupsRequest?: ObjectRef.AsObject
): Promise<UpstreamGroup.AsObject[]> {
  let request = new ListUpstreamGroupsRequest();
  if (listUpstreamGroupsRequest) {
    request.setGlooInstanceRef(
      getObjectRefClassFromRefObj(listUpstreamGroupsRequest)
    );
  }

  return new Promise((resolve, reject) => {
    glooResourceApiClient.listUpstreamGroups(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().upstreamGroupsList);
      }
    });
  });
}

function listProxies(
  listProxiesRequest?: ObjectRef.AsObject
): Promise<Proxy.AsObject[]> {
  let request = new ListProxiesRequest();
  if (listProxiesRequest) {
    request.setGlooInstanceRef(getObjectRefClassFromRefObj(listProxiesRequest));
  }

  return new Promise((resolve, reject) => {
    glooResourceApiClient.listProxies(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().proxiesList);
      }
    });
  });
}

function listSettings(
  listSettingsRequest?: ObjectRef.AsObject
): Promise<Settings.AsObject[]> {
  let request = new ListProxiesRequest();
  if (listSettingsRequest) {
    request.setGlooInstanceRef(
      getObjectRefClassFromRefObj(listSettingsRequest)
    );
  }

  return new Promise((resolve, reject) => {
    glooResourceApiClient.listSettings(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().settingsList);
      }
    });
  });
}

function getUpstream(
  glooInstRef: ObjectRef.AsObject,
  upstreamRef: ClusterObjectRef.AsObject
): Promise<Upstream.AsObject> {
  let request = new ListUpstreamsRequest();
  request.setGlooInstanceRef(getObjectRefClassFromRefObj(glooInstRef));

  return new Promise((resolve, reject) => {
    glooResourceApiClient.listUpstreams(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        const upstreamsList = data!.toObject().upstreamsList;
        const upstream = upstreamsList.find(
          u =>
            (!upstreamRef.clusterName ||
              u.metadata?.clusterName === upstreamRef.clusterName) &&
            u.metadata?.namespace === upstreamRef.namespace &&
            u.metadata?.name === upstreamRef.name
        );
        if (upstream) {
          resolve(upstream);
        } else {
          reject({ message: 'Upstream not found' });
        }
      }
    });
  });
}

function getUpstreamGroup(
  glooInstRef: ObjectRef.AsObject,
  upstreamGroupRef: ClusterObjectRef.AsObject
): Promise<UpstreamGroup.AsObject> {
  let request = new ListUpstreamGroupsRequest();
  request.setGlooInstanceRef(getObjectRefClassFromRefObj(glooInstRef));

  return new Promise((resolve, reject) => {
    glooResourceApiClient.listUpstreamGroups(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        const upstreamGroupsList = data!.toObject().upstreamGroupsList;
        const upstreamGroup = upstreamGroupsList.find(
          u =>
            (!upstreamGroupRef.clusterName ||
              u.metadata?.clusterName === upstreamGroupRef.clusterName) &&
            u.metadata?.namespace === upstreamGroupRef.namespace &&
            u.metadata?.name === upstreamGroupRef.name
        );
        if (upstreamGroup) {
          resolve(upstreamGroup);
        } else {
          reject({ message: 'Upstream Group not found' });
        }
      }
    });
  });
}

function getUpstreamYAML(
  upstreamObjectRef: ClusterObjectRef.AsObject
): Promise<string> {
  let request = new GetUpstreamYamlRequest();
  request.setUpstreamRef(
    getClusterRefClassFromClusterRefObj(upstreamObjectRef)
  );

  return new Promise((resolve, reject) => {
    glooResourceApiClient.getUpstreamYaml(request, (error, data) => {
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

function getUpstreamGroupYAML(
  upstreamGroupObjectRef: ClusterObjectRef.AsObject
): Promise<string> {
  let request = new GetUpstreamGroupYamlRequest();
  request.setUpstreamGroupRef(
    getClusterRefClassFromClusterRefObj(upstreamGroupObjectRef)
  );

  return new Promise((resolve, reject) => {
    glooResourceApiClient.getUpstreamGroupYaml(request, (error, data) => {
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

function getProxyYAML(
  proxyObjectRef: ClusterObjectRef.AsObject
): Promise<string> {
  let request = new GetProxyYamlRequest();
  request.setProxyRef(getClusterRefClassFromClusterRefObj(proxyObjectRef));

  return new Promise((resolve, reject) => {
    glooResourceApiClient.getProxyYaml(request, (error, data) => {
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

function getSettingYAML(
  settingObjectRef: ClusterObjectRef.AsObject
): Promise<string> {
  let request = new GetSettingsYamlRequest();
  request.setSettingsRef(getClusterRefClassFromClusterRefObj(settingObjectRef));

  return new Promise((resolve, reject) => {
    glooResourceApiClient.getSettingsYaml(request, (error, data) => {
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
