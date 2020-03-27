import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { grpc } from '@improbable-eng/grpc-web';
import { PortalApi } from 'proto/dev-portal/api/grpc/admin/portal_pb_service';
import { host } from 'store';
import {
  Portal,
  PortalList,
  PortalWriteRequest
} from 'proto/dev-portal/api/grpc/admin/portal_pb';
import {
  ObjectRef,
  Selector
} from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import {
  PortalSpec,
  PortalStatus,
  KeyScope
} from 'proto/dev-portal/api/dev-portal/v1/portal_pb';
import { Metadata } from 'proto/solo-kit/api/v1/metadata_pb';
import { ObjectMeta } from 'proto/dev-portal/api/grpc/common/common_pb';

export const devPortalApi = {
  listPortals,
  deletePortal,
  getPortal,
  createPortal
};

function listPortals(): Promise<Portal.AsObject[]> {
  return new Promise((resolve, reject) => {
    grpc.invoke(PortalApi.ListPortals, {
      request: new Empty(),
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: PortalList) => {
        if (message) {
          resolve(message.toObject().portalsList);
        }
      },
      onEnd: (
        status: grpc.Code,
        statusMessage: string,
        trailers: grpc.Metadata
      ) => {
        if (status !== grpc.Code.OK) {
          reject(statusMessage);
        }
      }
    });
  });
}

// service PortalApi {
//   // Returns a portal resource, without the corresponding static assets
//   rpc GetPortal (.devportal.solo.io.ObjectRef) returns (Portal) {
//   }
//   // Returns a portal resource, including the corresponding static assets
//   rpc GetPortalWithAssets (.devportal.solo.io.ObjectRef) returns (Portal) {
//   }
//   // Returns all portals (each without the corresponding static assets)
//   rpc ListPortals (google.protobuf.Empty) returns (PortalList) {
//   }
//   rpc CreatePortal (PortalWriteRequest) returns (Portal) {
//   }
//   rpc UpdatePortal (PortalWriteRequest) returns (Portal) {
//   }
//   rpc DeletePortal (.devportal.solo.io.ObjectRef) returns (google.protobuf.Empty) {
//   }
// }

function getPortal(portalRef: ObjectRef.AsObject): Promise<Portal.AsObject> {
  const { name, namespace } = portalRef;
  let request = new ObjectRef();
  request.setName(name);
  request.setNamespace(namespace);

  return new Promise((resolve, reject) => {
    grpc.invoke(PortalApi.GetPortal, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: Portal) => {
        if (message) {
          resolve(message.toObject());
        }
      },
      onEnd: (
        status: grpc.Code,
        statusMessage: string,
        trailers: grpc.Metadata
      ) => {
        if (status !== grpc.Code.OK) {
          reject(statusMessage);
        }
      }
    });
  });
}

function getPortalWithAssets(
  portalRef: ObjectRef.AsObject
): Promise<Portal.AsObject> {
  const { name, namespace } = portalRef;
  let request = new ObjectRef();
  request.setName(name);
  request.setNamespace(namespace);

  return new Promise((resolve, reject) => {
    grpc.invoke(PortalApi.GetPortal, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: Portal) => {
        if (message) {
          resolve(message.toObject());
        }
      },
      onEnd: (
        status: grpc.Code,
        statusMessage: string,
        trailers: grpc.Metadata
      ) => {
        if (status !== grpc.Code.OK) {
          reject(statusMessage);
        }
      }
    });
  });
}

function setPortalValuesToGrpc(
  portal: Portal.AsObject,
  portalToUpdate = new Portal()
): Portal {
  let { spec, metadata, status } = portal!;
  if (metadata !== undefined) {
    let { name, namespace } = metadata;
    let newMetadata = new ObjectMeta();
    newMetadata.setName(name);
    newMetadata.setNamespace(namespace);
    portalToUpdate.setMetadata(newMetadata);
  }

  if (spec !== undefined) {
    let newSpec = new PortalSpec();
    const {
      description,
      displayName,
      domainsList,
      keyScopesList, //
      staticPagesList, //
      banner, //
      customStyling, //
      favicon, //
      primaryLogo, //
      publishApiDocs //
    } = spec;
    if (description !== undefined) {
      newSpec.setDescription(description);
    }
    if (displayName !== undefined) {
      newSpec.setDisplayName(displayName);
    }

    if (domainsList !== undefined) {
      newSpec.setDomainsList(domainsList);
    }

    if (keyScopesList !== undefined) {
      let newKeyScopesList = keyScopesList.map(keyScopeObj => {
        const { name, namespace, description, apiDocs } = keyScopeObj;

        let keyScope = new KeyScope();
        keyScope.setName(name);
        keyScope.setNamespace(namespace);
        keyScope.setDescription(description);
        let matchLabelsMapSelector = new Selector();

        apiDocs?.matchLabelsMap.forEach(([key, value], idx) =>
          matchLabelsMapSelector.getMatchLabelsMap().set(key, value)
        );
        keyScope.setApiDocs(matchLabelsMapSelector);
      });
    }

    portalToUpdate.setSpec();
  }

  return portalToUpdate;
}

function createPortal(
  portalWriteRequest: PortalWriteRequest.AsObject
): Promise<Portal.AsObject> {
  const { portal, usersList, apiDocsList, groupsList } = portalWriteRequest;
  let request = new PortalWriteRequest();
  let portalToCreate = new Portal();

  if (portal !== undefined) {
    const { spec, status, metadata } = portal;
    let portalSpecToCreate = new PortalSpec();
    let portalStatusToCreate = new PortalStatus();
  }
  if (apiDocsList !== undefined) {
    let apiDocsRefList = apiDocsList.map(apiDocRefObj => {
      let apiDocRef = new ObjectRef();
      apiDocRef.setName(apiDocRefObj.name);
      apiDocRef.setNamespace(apiDocRefObj.namespace);
      return apiDocRef;
    });
    request.setApiDocsList(apiDocsRefList);
  }

  if (usersList !== undefined) {
    let usersRefList = usersList.map(userRefObj => {
      let userRef = new ObjectRef();
      userRef.setName(userRefObj.name);
      userRef.setNamespace(userRefObj.namespace);
      return userRef;
    });
    request.setUsersList(usersRefList);
  }

  if (groupsList !== undefined) {
    let groupsRefList = groupsList.map(groupRefObj => {
      let groupRef = new ObjectRef();
      groupRef.setName(groupRefObj.name);
      groupRef.setNamespace(groupRefObj.namespace);
      return groupRef;
    });
    request.setGroupsList(groupsRefList);
  }

  return new Promise((resolve, reject) => {
    grpc.invoke(PortalApi.CreatePortal, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: Portal) => {
        if (message) {
          resolve(message.toObject());
        }
      },
      onEnd: (
        status: grpc.Code,
        statusMessage: string,
        trailers: grpc.Metadata
      ) => {
        if (status !== grpc.Code.OK) {
          reject(statusMessage);
        }
      }
    });
  });
}

function deletePortal(portalRef: ObjectRef.AsObject): Promise<Empty.AsObject> {
  const { name, namespace } = portalRef;
  let request = new ObjectRef();
  request.setName(name);
  request.setNamespace(namespace);

  return new Promise((resolve, reject) => {
    grpc.invoke(PortalApi.DeletePortal, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: Empty) => {
        if (message) {
          resolve(message.toObject());
        }
      },
      onEnd: (
        status: grpc.Code,
        statusMessage: string,
        trailers: grpc.Metadata
      ) => {
        if (status !== grpc.Code.OK) {
          reject(statusMessage);
        }
      }
    });
  });
}
