import { grpc } from '@improbable-eng/grpc-web';
import {
  ClusterObjectRef,
  ObjectRef,
} from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import {
  Pagination,
  StatusFilter,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb';
import {
  GetProxyDetailsRequest,
  GetProxyDetailsResponse,
  GetProxyYamlRequest,
  GetSettingsDetailsRequest,
  GetSettingsDetailsResponse,
  GetSettingsYamlRequest,
  GetUpstreamDetailsRequest,
  GetUpstreamDetailsResponse,
  GetUpstreamGroupDetailsRequest,
  GetUpstreamGroupDetailsResponse,
  GetUpstreamGroupYamlRequest,
  GetUpstreamYamlRequest,
  ListProxiesRequest,
  ListUpstreamGroupsRequest,
  ListUpstreamGroupsResponse,
  ListUpstreamsRequest,
  ListUpstreamsResponse,
  Proxy,
  Settings,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import { GlooResourceApiClient } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb_service';
import {
  getClusterRefClassFromClusterRefObj,
  getObjectRefClassFromRefObj,
  host,
  toPaginationClass,
} from './helpers';

const glooResourceApiClient = new GlooResourceApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true,
});

export const glooResourceApi = {
  listUpstreams,
  listUpstreamGroups,
  listProxies,
  listSettings,
  getUpstreamYAML,
  getUpstreamGroupYAML,
  getProxyYAML,
  getSettingsYAML,
  getUpstreamDetails,
  getUpstreamGroupDetails,
  getProxyDetails,
  getSettingsDetails,
};

function listUpstreams(
  listUpstreamsRequest?: ObjectRef.AsObject,
  pagination?: Pagination.AsObject,
  queryString?: string,
  statusFilter?: number
): Promise<ListUpstreamsResponse.AsObject> {
  // Used to debug slowdown issues.
  // return new Promise((resolve, reject) => {
  //   setTimeout(() => {
  //     const res = new ListUpstreamsResponse();
  //     res.setTotal(1903);
  //     const ul = [] as Upstream[];
  //     for (let i = 0; i < res.getTotal(); i++) {
  //       const u = new Upstream();
  //       const g = new ClusterObjectRef();
  //       g.setName('gloo');
  //       g.setNamespace('gloo-system');
  //       u.setGlooInstance(g);
  //       const m = new ObjectMeta();
  //       m.setName('upstream_' + i);
  //       m.setNamespace('gloo-system');
  //       m.setUid(i.toString());
  //       u.setMetadata(m);
  //       ul.push(u);
  //     }
  //     res.setUpstreamsList(ul);
  //     resolve(res!.toObject());
  //   }, 500);
  // });
  let request = new ListUpstreamsRequest();
  if (listUpstreamsRequest) {
    request.setGlooInstanceRef(
      getObjectRefClassFromRefObj(listUpstreamsRequest)
    );
  }
  if (pagination) request.setPagination(toPaginationClass(pagination));
  if (queryString) request.setQueryString(queryString);
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
  listUpstreamGroupsRequest?: ObjectRef.AsObject,
  pagination?: Pagination.AsObject,
  queryString?: string,
  statusFilter?: number
): Promise<ListUpstreamGroupsResponse.AsObject> {
  let request = new ListUpstreamGroupsRequest();
  if (listUpstreamGroupsRequest) {
    request.setGlooInstanceRef(
      getObjectRefClassFromRefObj(listUpstreamGroupsRequest)
    );
  }
  if (pagination) request.setPagination(toPaginationClass(pagination));
  if (queryString) request.setQueryString(queryString);
  if (statusFilter !== undefined) {
    const sf = new StatusFilter();
    sf.setState(statusFilter);
    request.setStatusFilter(sf);
  }
  return new Promise((resolve, reject) => {
    glooResourceApiClient.listUpstreamGroups(request, (error, data) => {
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

function getSettingsYAML(
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

export function getUpstreamDetails(
  upstreamObjectRef: ClusterObjectRef.AsObject
): Promise<GetUpstreamDetailsResponse.AsObject> {
  let request = new GetUpstreamDetailsRequest();
  request.setUpstreamRef(
    getClusterRefClassFromClusterRefObj(upstreamObjectRef)
  );
  return new Promise((resolve, reject) => {
    glooResourceApiClient.getUpstreamDetails(request, (error, data) => {
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

function getUpstreamGroupDetails(
  upstreamGroupObjectRef: ClusterObjectRef.AsObject
): Promise<GetUpstreamGroupDetailsResponse.AsObject> {
  let request = new GetUpstreamGroupDetailsRequest();
  request.setUpstreamGroupRef(
    getClusterRefClassFromClusterRefObj(upstreamGroupObjectRef)
  );
  return new Promise((resolve, reject) => {
    glooResourceApiClient.getUpstreamGroupDetails(request, (error, data) => {
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

function getProxyDetails(
  proxyObjectRef: ClusterObjectRef.AsObject
): Promise<GetProxyDetailsResponse.AsObject> {
  let request = new GetProxyDetailsRequest();
  request.setProxyRef(getClusterRefClassFromClusterRefObj(proxyObjectRef));
  return new Promise((resolve, reject) => {
    glooResourceApiClient.getProxyDetails(request, (error, data) => {
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

function getSettingsDetails(
  settingObjectRef: ClusterObjectRef.AsObject
): Promise<GetSettingsDetailsResponse.AsObject> {
  let request = new GetSettingsDetailsRequest();
  request.setSettingsRef(getClusterRefClassFromClusterRefObj(settingObjectRef));
  return new Promise((resolve, reject) => {
    glooResourceApiClient.getSettingsDetails(request, (error, data) => {
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
