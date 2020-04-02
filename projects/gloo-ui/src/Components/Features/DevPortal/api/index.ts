import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { grpc } from '@improbable-eng/grpc-web';
import { PortalApi } from 'proto/dev-portal/api/grpc/admin/portal_pb_service';
import { ApiDocApi } from 'proto/dev-portal/api/grpc/admin/apidoc_pb_service';
import { UserApi } from 'proto/dev-portal/api/grpc/admin/user_pb_service';
import { GroupApi } from 'proto/dev-portal/api/grpc/admin/group_pb_service';
import { ApiKeyApi } from 'proto/dev-portal/api/grpc/admin/api_key_pb_service';
import { ApiKeyScopeApi } from 'proto/dev-portal/api/grpc/admin/api_key_scope_pb_service';

import { host } from 'store';

import {
  ApiDoc,
  ApiDocList,
  ApiDocFilter,
  ApiDocGetRequest,
  ApiDocWriteRequest
} from 'proto/dev-portal/api/grpc/admin/apidoc_pb';
import {
  Portal,
  PortalList,
  PortalWriteRequest
} from 'proto/dev-portal/api/grpc/admin/portal_pb';

import {
  User,
  UserList,
  UserFilter,
  UserWriteRequest
} from 'proto/dev-portal/api/grpc/admin/user_pb';
import {
  Group,
  GroupList,
  GroupFilter,
  GroupWriteRequest
} from 'proto/dev-portal/api/grpc/admin/group_pb';
import { ApiKey, ApiKeyList } from 'proto/dev-portal/api/grpc/admin/api_key_pb';

import {
  ObjectRef,
  Selector,
  DataSource
} from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import {
  PortalSpec,
  PortalStatus,
  KeyScope,
  KeyScopeStatus,
  CustomStyling,
  StaticPage
} from 'proto/dev-portal/api/dev-portal/v1/portal_pb';
import { ObjectMeta, Time } from 'proto/dev-portal/api/grpc/common/common_pb';
import {
  createDataSourceClassFromObject,
  createPortalClassFromObject
} from '../api-helper';
import {
  ApiDocStatus,
  ApiDocSpec
} from 'proto/dev-portal/api/dev-portal/v1/apidoc_pb';
import {
  UserStatus,
  UserSpec
} from 'proto/dev-portal/api/dev-portal/v1/user_pb';
import {
  AccessLevel,
  AccessLevelStatus
} from 'proto/dev-portal/api/dev-portal/v1/access_level_pb';
import {
  GroupStatus,
  GroupSpec
} from 'proto/dev-portal/api/dev-portal/v1/group_pb';
import {
  ApiKeyScopeWithApiDocs,
  ApiKeyScopeList,
  ApiKeyScopeRef,
  ApiKeyScopeWriteRequest,
  ApiKeyScope
} from 'proto/dev-portal/api/grpc/admin/api_key_scope_pb';

export const portalApi = {
  listPortals,
  deletePortal,
  createPortal,
  updatePortal,
  getPortalWithAssets,
  createPortalPage,
  updatePortalPage,
  deletePortalPage
};

export const apiDocApi = {
  listApiDocs,
  getApiDoc,
  createApiDoc,
  deleteApiDoc
};

export const userApi = {
  listUsers,
  createUser
};

export const groupApi = {
  listGroups,
  createGroup
};

export const apiKeyApi = {
  listApiKeys,
  deleteApiKey
};

export const apiKeyScopeApi = {
  listKeyScopes,
  createKeyScope,
  updateKeyScope,
  deleteKeyScope
};

function deleteApiDoc(apiDocRef: ObjectRef.AsObject): Promise<Empty.AsObject> {
  const { name, namespace } = apiDocRef;
  let request = new ObjectRef();
  request.setName(name);
  request.setNamespace(namespace);

  return new Promise((resolve, reject) => {
    grpc.invoke(ApiDocApi.DeleteApiDoc, {
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

export function groupMessageFromObject(
  group: Group.AsObject,
  groupToUpdate = new Group()
): Group {
  let { spec, metadata, status } = group!;
  if (metadata !== undefined) {
    let { name, namespace } = metadata;
    let newMetadata = new ObjectMeta();
    newMetadata.setName(name);
    newMetadata.setNamespace(namespace);
    groupToUpdate.setMetadata(newMetadata);
  }

  if (status !== undefined) {
    let statusMessage = new GroupStatus();
    const {
      usersList,
      accessLevel,
      observedGeneration,
      reason,
      state
    } = status;

    if (usersList !== undefined) {
      let userRefList = usersList.map(userRefObj => {
        let userRef = new ObjectRef();
        userRef.setName(userRefObj.name);
        userRef.setNamespace(userRefObj.namespace);
        return userRef;
      });
      statusMessage.setUsersList(userRefList);
    }
    if (accessLevel !== undefined) {
      let newAccessLevel = new AccessLevelStatus();
      const { apiDocsList, portalsList } = accessLevel;
      if (apiDocsList !== undefined) {
        let apiDocRefList = apiDocsList.map(apiDocObj => {
          let apiDocRef = new ObjectRef();
          apiDocRef.setName(apiDocObj.name);
          apiDocRef.setNamespace(apiDocObj.namespace);
          return apiDocRef;
        });
        newAccessLevel.setApiDocsList(apiDocRefList);
      }
      if (portalsList !== undefined) {
        let portalRefList = portalsList.map(portalRefObj => {
          let portalRef = new ObjectRef();
          portalRef.setName(portalRefObj.name);
          portalRef.setNamespace(portalRefObj.namespace);
          return portalRef;
        });
        newAccessLevel.setPortalsList(portalRefList);
      }
      statusMessage.setAccessLevel(newAccessLevel);
    }

    if (observedGeneration !== undefined) {
      statusMessage.setObservedGeneration(observedGeneration);
    }

    if (reason !== undefined) {
      statusMessage.setReason(reason);
    }
    if (state !== undefined) {
      statusMessage.setState(state);
    }
    groupToUpdate.setStatus(statusMessage);
  }

  if (spec !== undefined) {
    let newSpec = new GroupSpec();
    const { displayName, description, userSelector, accessLevel } = spec;

    if (displayName !== undefined) {
      newSpec.setDisplayName(displayName);
    }
    if (description !== undefined) {
      newSpec.setDescription(description);
    }

    if (userSelector !== undefined) {
      let userSelectorMessage = selectorMessageFromObject(userSelector);
      newSpec.setUserSelector(userSelectorMessage);
    }
    if (accessLevel !== undefined) {
      let newAccessLevel = new AccessLevel();
      const { apiDocSelector, portalSelector } = accessLevel;
      if (apiDocSelector !== undefined) {
        let apiDocSelectorMessage = selectorMessageFromObject(apiDocSelector);
        newAccessLevel.setApiDocSelector(apiDocSelectorMessage);
      }

      if (portalSelector !== undefined) {
        let portalSelectorMessage = selectorMessageFromObject(portalSelector);
        newAccessLevel.setApiDocSelector(portalSelectorMessage);
      }

      newSpec.setAccessLevel(newAccessLevel);
    }

    groupToUpdate.setSpec(newSpec);
  }

  return groupToUpdate;
}

function createGroup(
  groupWriteRequest: GroupWriteRequest.AsObject
): Promise<Group.AsObject> {
  const { usersList, group, apiDocsList, portalsList } = groupWriteRequest;
  let request = new GroupWriteRequest();

  if (group !== undefined) {
    let groupToCreate = groupMessageFromObject(group);
    request.setGroup(groupToCreate);
  }

  if (portalsList !== undefined) {
    let portalRefList = portalsList.map(portalRefObj => {
      let portalRef = new ObjectRef();
      portalRef.setName(portalRefObj.name);
      portalRef.setNamespace(portalRefObj.namespace);
      return portalRef;
    });
    request.setPortalsList(portalRefList);
  }

  if (apiDocsList !== undefined) {
    let apiDocRefList = apiDocsList.map(apiDocObj => {
      let apiDocRef = new ObjectRef();
      apiDocRef.setName(apiDocObj.name);
      apiDocRef.setNamespace(apiDocObj.namespace);
      return apiDocRef;
    });
    request.setApiDocsList(apiDocRefList);
  }

  if (usersList !== undefined) {
    let userRefList = usersList.map(userRefObj => {
      let userRef = new ObjectRef();
      userRef.setName(userRefObj.name);
      userRef.setNamespace(userRefObj.namespace);
      return userRef;
    });
    request.setUsersList(userRefList);
  }

  return new Promise((resolve, reject) => {
    grpc.invoke(GroupApi.CreateGroup, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: Group) => {
        if (message) {
          resolve(message.toObject());
        }
      },
      onEnd: (
        status: grpc.Code,
        statusMessage: string,
        trailers: grpc.Metadata
      ) => {
        console.log('statusMessage', statusMessage);
        if (status !== grpc.Code.OK) {
          reject(statusMessage);
        }
      }
    });
  });
}

export function userMessageFromObject(
  user: User.AsObject,
  userToUpdate = new User()
): User {
  let { spec, metadata, status } = user!;
  if (metadata !== undefined) {
    let { name, namespace } = metadata;
    let newMetadata = new ObjectMeta();
    newMetadata.setName(name);
    newMetadata.setNamespace(namespace);
    userToUpdate.setMetadata(newMetadata);
  }

  if (status !== undefined) {
    let statusMessage = new UserStatus();
    const {
      hasLoggedIn,
      accessLevel,
      observedGeneration,
      reason,
      state
    } = status;

    if (hasLoggedIn !== undefined) {
      statusMessage.setHasLoggedIn(hasLoggedIn);
    }
    if (accessLevel !== undefined) {
      let newAccessLevel = new AccessLevelStatus();
      const { apiDocsList, portalsList } = accessLevel;
      if (apiDocsList !== undefined) {
        let apiDocRefList = apiDocsList.map(apiDocObj => {
          let apiDocRef = new ObjectRef();
          apiDocRef.setName(apiDocObj.name);
          apiDocRef.setNamespace(apiDocObj.namespace);
          return apiDocRef;
        });
        newAccessLevel.setApiDocsList(apiDocRefList);
      }
      if (portalsList !== undefined) {
        let portalRefList = portalsList.map(portalRefObj => {
          let portalRef = new ObjectRef();
          portalRef.setName(portalRefObj.name);
          portalRef.setNamespace(portalRefObj.namespace);
          return portalRef;
        });
        newAccessLevel.setPortalsList(portalRefList);
      }
      statusMessage.setAccessLevel(newAccessLevel);
    }

    if (observedGeneration !== undefined) {
      statusMessage.setObservedGeneration(observedGeneration);
    }

    if (reason !== undefined) {
      statusMessage.setReason(reason);
    }
    if (state !== undefined) {
      statusMessage.setState(state);
    }
    userToUpdate.setStatus(statusMessage);
  }

  if (spec !== undefined) {
    let newSpec = new UserSpec();
    const { email, username, accessLevel, basicAuth } = spec;

    if (email !== undefined) {
      newSpec.setEmail(email);
    }
    if (username !== undefined) {
      newSpec.setUsername(username);
    }

    if (accessLevel !== undefined) {
      let newAccessLevel = new AccessLevel();
      const { apiDocSelector, portalSelector } = accessLevel;
      if (apiDocSelector !== undefined) {
        let apiDocSelectorMessage = selectorMessageFromObject(apiDocSelector);
        newAccessLevel.setApiDocSelector(apiDocSelectorMessage);
      }

      if (portalSelector !== undefined) {
        let portalSelectorMessage = selectorMessageFromObject(portalSelector);
        newAccessLevel.setApiDocSelector(portalSelectorMessage);
      }

      newSpec.setAccessLevel(newAccessLevel);
    }

    if (basicAuth !== undefined) {
      let newBasicAuth = new UserSpec.BasicAuth();
      const {
        passwordSecretKey,
        passwordSecretName,
        passwordSecretNamespace
      } = basicAuth;
      if (passwordSecretKey !== undefined) {
        newBasicAuth.setPasswordSecretKey(passwordSecretKey);
      }
      if (passwordSecretName !== undefined) {
        newBasicAuth.setPasswordSecretName(passwordSecretName);
      }
      if (passwordSecretNamespace !== undefined) {
        newBasicAuth.setPasswordSecretNamespace(passwordSecretNamespace);
      }

      newSpec.setBasicAuth(newBasicAuth);
    }

    userToUpdate.setSpec(newSpec);
  }

  return userToUpdate;
}

function createUser(
  userWriteRequest: UserWriteRequest.AsObject
): Promise<User.AsObject> {
  const {
    user,
    apiDocsList,
    password,
    groupsList,
    portalsList
  } = userWriteRequest;
  let request = new UserWriteRequest();

  if (user !== undefined) {
    let userToCreate = userMessageFromObject(user);
    request.setUser(userToCreate);
  }
  if (password !== undefined) {
    request.setPassword(password);
  }

  if (portalsList !== undefined) {
    let portalRefList = portalsList.map(portalRefObj => {
      let portalRef = new ObjectRef();
      portalRef.setName(portalRefObj.name);
      portalRef.setNamespace(portalRefObj.namespace);
      return portalRef;
    });
    request.setPortalsList(portalRefList);
  }

  if (apiDocsList !== undefined) {
    let apiDocRefList = apiDocsList.map(apiDocObj => {
      let apiDocRef = new ObjectRef();
      apiDocRef.setName(apiDocObj.name);
      apiDocRef.setNamespace(apiDocObj.namespace);
      return apiDocRef;
    });
    request.setApiDocsList(apiDocRefList);
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
    grpc.invoke(UserApi.CreateUser, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: User) => {
        if (message) {
          resolve(message.toObject());
        }
      },
      onEnd: (
        status: grpc.Code,
        statusMessage: string,
        trailers: grpc.Metadata
      ) => {
        console.log('statusMessage', statusMessage);
        if (status !== grpc.Code.OK) {
          reject(statusMessage);
        }
      }
    });
  });
}

export function apiDocMessageFromObject(
  apiDoc: ApiDoc.AsObject,
  apiDocToUpdate = new ApiDoc()
): ApiDoc {
  let { spec, metadata, status } = apiDoc!;
  if (metadata !== undefined) {
    let { name, namespace } = metadata;
    let newMetadata = new ObjectMeta();
    newMetadata.setName(name);
    newMetadata.setNamespace(namespace);
    apiDocToUpdate.setMetadata(newMetadata);
  }

  if (status !== undefined) {
    let statusMessage = new ApiDocStatus();
    const {
      basePath,
      description,
      displayName,
      numberOfEndpoints,
      version,
      modifiedDate,
      observedGeneration,
      reason,
      state
    } = status;
    if (basePath !== undefined) {
      statusMessage.setBasePath(basePath);
    }
    if (description !== undefined) {
      statusMessage.setDescription(description);
    }
    if (displayName !== undefined) {
      statusMessage.setDisplayName(displayName);
    }

    if (numberOfEndpoints !== undefined) {
      statusMessage.setNumberOfEndpoints(numberOfEndpoints);
    }

    if (version !== undefined) {
      statusMessage.setVersion(version);
    }

    if (observedGeneration !== undefined) {
      statusMessage.setObservedGeneration(observedGeneration);
    }

    if (reason !== undefined) {
      statusMessage.setReason(reason);
    }
    if (state !== undefined) {
      statusMessage.setState(state);
    }
    apiDocToUpdate.setStatus(statusMessage);
  }

  if (spec !== undefined) {
    let newSpec = new ApiDocSpec();
    const { dataSource, image, openApi } = spec;

    if (dataSource !== undefined) {
      let dataSourceMessage = dataSourceMessageFromObject(dataSource);
      newSpec.setDataSource(dataSourceMessage);
    }

    if (image !== undefined) {
      let imageMessage = dataSourceMessageFromObject(image);
      newSpec.setImage(imageMessage);
    }

    if (openApi !== undefined) {
      let openApiMessage = new ApiDocSpec.OpenApi();

      newSpec.setOpenApi(openApiMessage);
    }

    apiDocToUpdate.setSpec(newSpec);
  }

  return apiDocToUpdate;
}

function createApiDoc(
  apiDocWriteRequest: ApiDocWriteRequest.AsObject
): Promise<ApiDoc.AsObject> {
  const { apidoc, usersList, groupsList, portalsList } = apiDocWriteRequest;
  let request = new ApiDocWriteRequest();
  if (apidoc !== undefined) {
    let apiDocToCreate = apiDocMessageFromObject(apidoc);
    request.setApidoc(apiDocToCreate);
  }

  if (portalsList !== undefined) {
    let portalRefList = portalsList.map(portalRefObj => {
      let portalRef = new ObjectRef();
      portalRef.setName(portalRefObj.name);
      portalRef.setNamespace(portalRefObj.namespace);
      return portalRef;
    });
    request.setPortalsList(portalRefList);
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
    grpc.invoke(ApiDocApi.CreateApiDoc, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: ApiDoc) => {
        if (message) {
          resolve(message.toObject());
        }
      },
      onEnd: (
        status: grpc.Code,
        statusMessage: string,
        trailers: grpc.Metadata
      ) => {
        console.log('statusMessage', statusMessage);
        if (status !== grpc.Code.OK) {
          reject(statusMessage);
        }
      }
    });
  });
}

function listApiDocs(
  apiDocFilter: ApiDocFilter.AsObject
): Promise<ApiDoc.AsObject[]> {
  const { portalsList } = apiDocFilter;
  let request = new ApiDocFilter();

  if (portalsList !== undefined) {
    let portalsRefList = portalsList.map(portalObj => {
      let portalRef = new ObjectRef();
      portalRef.setName(portalObj.name);
      portalRef.setNamespace(portalObj.namespace);
      return portalRef;
    });
    request.setPortalsList(portalsRefList);
  }
  return new Promise((resolve, reject) => {
    grpc.invoke(ApiDocApi.ListApiDocs, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: ApiDocList) => {
        if (message) {
          resolve(message.toObject().apidocsList);
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

function getApiDoc(
  getApiDocRequest: ApiDocGetRequest.AsObject
): Promise<ApiDoc.AsObject> {
  const { apidoc, withassets } = getApiDocRequest;
  let request = new ApiDocGetRequest();
  let apiDocRef = new ObjectRef();
  if (apidoc !== undefined) {
    apiDocRef.setName(apidoc.name);
    apiDocRef.setNamespace(apidoc.namespace);
    request.setApidoc(apiDocRef);
  }
  if (withassets !== undefined) {
    request.setWithassets(withassets);
  }

  return new Promise((resolve, reject) => {
    grpc.invoke(ApiDocApi.GetApiDoc, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: ApiDoc) => {
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

function deleteApiKey(apiKeyRef: ObjectRef.AsObject): Promise<Empty.AsObject> {
  const { name, namespace } = apiKeyRef;
  let request = new ObjectRef();
  request.setName(name);
  request.setNamespace(namespace);

  return new Promise((resolve, reject) => {
    grpc.invoke(ApiKeyApi.DeleteApiKey, {
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

function listApiKeys(): Promise<ApiKey.AsObject[]> {
  return new Promise((resolve, reject) => {
    grpc.invoke(ApiKeyApi.ListApiKeys, {
      request: new Empty(),
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: ApiKeyList) => {
        if (message) {
          resolve(message.toObject().apiKeysList);
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

function listGroups(
  groupFilter: GroupFilter.AsObject
): Promise<Group.AsObject[]> {
  const { portalsList } = groupFilter;
  let request = new GroupFilter();

  if (portalsList !== undefined) {
    let portalsRefList = portalsList.map(portalObj => {
      let portalRef = new ObjectRef();
      portalRef.setName(portalObj.name);
      portalRef.setNamespace(portalObj.namespace);
      return portalRef;
    });
    request.setPortalsList(portalsRefList);
  }
  return new Promise((resolve, reject) => {
    grpc.invoke(GroupApi.ListGroups, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: GroupList) => {
        if (message) {
          resolve(message.toObject().groupsList);
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

function listUsers(userFilter: UserFilter.AsObject): Promise<User.AsObject[]> {
  const { portalsList } = userFilter;
  let request = new UserFilter();

  if (portalsList !== undefined) {
    let portalsRefList = portalsList.map(portalObj => {
      let portalRef = new ObjectRef();
      portalRef.setName(portalObj.name);
      portalRef.setNamespace(portalObj.namespace);
      return portalRef;
    });
    request.setPortalsList(portalsRefList);
  }
  return new Promise((resolve, reject) => {
    grpc.invoke(UserApi.ListUsers, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: UserList) => {
        if (message) {
          resolve(message.toObject().usersList);
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

function getPortalWithAssets(
  portalRef: ObjectRef.AsObject
): Promise<Portal.AsObject> {
  const { name, namespace } = portalRef;
  let request = new ObjectRef();
  request.setName(name);
  request.setNamespace(namespace);

  return new Promise((resolve, reject) => {
    grpc.invoke(PortalApi.GetPortalWithAssets, {
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

export function portalMessageFromObject(
  portal: Portal.AsObject,
  portalToUpdate = new Portal()
): Portal {
  let { spec, metadata, status } = portal!;
  if (metadata !== undefined) {
    let { name, namespace, resourceVersion } = metadata;
    let newMetadata = new ObjectMeta();
    newMetadata.setName(name);
    newMetadata.setNamespace(namespace);
    newMetadata.setResourceVersion(resourceVersion);
    portalToUpdate.setMetadata(newMetadata);
  }

  if (status !== undefined) {
    let statusMessage = new PortalStatus();
    const {
      apiDocsList,
      keyScopesList,
      observedGeneration,
      publishUrl,
      reason,
      state
    } = status;
    if (apiDocsList !== undefined) {
      let apiDocsRefList = apiDocsList.map(apiDocObj => {
        let apiDocRef = new ObjectRef();
        apiDocRef.setName(apiDocObj.name);
        apiDocRef.setNamespace(apiDocObj.namespace);
        return apiDocRef;
      });
      statusMessage.setApiDocsList(apiDocsRefList);
    }

    if (keyScopesList !== undefined) {
      let keyScopeStatusList = keyScopesList.map(keyScopeStatusObj => {
        const {
          accessibleApiDocsList,
          name,
          provisionedKeysList
        } = keyScopeStatusObj;
        let keyScopeStatus = new KeyScopeStatus();
        if (name !== undefined) {
          keyScopeStatus.setName(name);
        }
        if (accessibleApiDocsList !== undefined) {
          let accessibleApiDocsListRefs = accessibleApiDocsList.map(
            accessibleApiDocObj => {
              let accessibleApiDocRef = new ObjectRef();
              accessibleApiDocRef.setName(accessibleApiDocObj.name);
              accessibleApiDocRef.setNamespace(accessibleApiDocObj.namespace);
              return accessibleApiDocRef;
            }
          );
          keyScopeStatus.setAccessibleApiDocsList(accessibleApiDocsListRefs);
        }

        if (provisionedKeysList !== undefined) {
          let provisionedKeysRefList = provisionedKeysList.map(
            provisionedKeyObj => {
              let provisionedKeyRef = new ObjectRef();
              provisionedKeyRef.setName(provisionedKeyObj.name);
              provisionedKeyRef.setNamespace(provisionedKeyObj.namespace);
              return provisionedKeyRef;
            }
          );
          keyScopeStatus.setProvisionedKeysList(provisionedKeysRefList);
        }
        if (accessibleApiDocsList !== undefined) {
          let accessibleApiDocsRefList = accessibleApiDocsList.map(
            accessibleApiDocObj => {
              let accessibleApiDocsRef = new ObjectRef();
              accessibleApiDocsRef.setName(accessibleApiDocObj.name);
              accessibleApiDocsRef.setNamespace(accessibleApiDocObj.namespace);
              return accessibleApiDocsRef;
            }
          );
          keyScopeStatus.setAccessibleApiDocsList(accessibleApiDocsRefList);
        }
        return keyScopeStatus;
      });
      statusMessage.setKeyScopesList(keyScopeStatusList);
    }

    if (observedGeneration !== undefined) {
      statusMessage.setObservedGeneration(observedGeneration);
    }
    if (publishUrl !== undefined) {
      statusMessage.setPublishUrl(publishUrl);
    }
    if (reason !== undefined) {
      statusMessage.setReason(reason);
    }
    if (state !== undefined) {
      statusMessage.setState(state);
    }
    portalToUpdate.setStatus(statusMessage);
  }

  if (spec !== undefined) {
    let newSpec = new PortalSpec();
    const {
      description,
      displayName,
      domainsList,
      keyScopesList,
      staticPagesList,
      banner,
      customStyling,
      favicon,
      primaryLogo,
      publishApiDocs
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
        return keyScope;
      });
      newSpec.setKeyScopesList(newKeyScopesList);
    }

    if (staticPagesList !== undefined) {
      let staticPagesMessageList = staticPagesList.map(staticPageObj => {
        const {
          description,
          name,
          navigationLinkName,
          path,
          content
        } = staticPageObj;
        let staticPageMessage = new StaticPage();
        if (description !== undefined) {
          staticPageMessage.setDescription(description);
        }
        if (name !== undefined) {
          staticPageMessage.setName(name);
        }
        if (navigationLinkName !== undefined) {
          staticPageMessage.setNavigationLinkName(navigationLinkName);
        }
        if (path !== undefined) {
          staticPageMessage.setPath(path);
        }

        if (content !== undefined) {
          let newContent = dataSourceMessageFromObject(content);

          staticPageMessage.setContent(newContent);
        }
        return staticPageMessage;
      });

      newSpec.setStaticPagesList(staticPagesMessageList);
    }

    if (banner !== undefined) {
      let bannerMessage = dataSourceMessageFromObject(banner);

      newSpec.setBanner(bannerMessage);
    }
    if (customStyling !== undefined) {
      let customStylingMessage = new CustomStyling();
      const {
        backgroundColor,
        buttonColorOverride,
        defaultTextColor,
        navigationLinksColorOverride,
        primaryColor,
        secondaryColor
      } = customStyling;
      if (backgroundColor !== undefined) {
        customStylingMessage.setBackgroundColor(backgroundColor);
      }
      if (buttonColorOverride !== undefined) {
        customStylingMessage.setButtonColorOverride(buttonColorOverride);
      }
      if (defaultTextColor !== undefined) {
        customStylingMessage.setDefaultTextColor(defaultTextColor);
      }
      if (navigationLinksColorOverride !== undefined) {
        customStylingMessage.setNavigationLinksColorOverride(
          navigationLinksColorOverride
        );
      }
      if (primaryColor !== undefined) {
        customStylingMessage.setPrimaryColor(primaryColor);
      }
      if (secondaryColor !== undefined) {
        customStylingMessage.setSecondaryColor(secondaryColor);
      }

      newSpec.setCustomStyling(customStylingMessage);
    }

    if (favicon !== undefined) {
      let faviconMessage = dataSourceMessageFromObject(favicon);
      newSpec.setFavicon(faviconMessage);
    }
    if (primaryLogo !== undefined) {
      let primaryLogoMessage = dataSourceMessageFromObject(primaryLogo);

      newSpec.setPrimaryLogo(primaryLogoMessage);
    }

    if (publishApiDocs !== undefined) {
      let publishApiDocsMessage = selectorMessageFromObject(publishApiDocs);
      newSpec.setPublishApiDocs(publishApiDocsMessage);
    }

    portalToUpdate.setSpec(newSpec);
  }

  return portalToUpdate;
}

function selectorMessageFromObject(
  selectorObj: Selector.AsObject,
  selectorMessage = new Selector()
): Selector {
  if (selectorObj.matchLabelsMap !== undefined) {
    selectorObj.matchLabelsMap.forEach(([key, value], idx) =>
      selectorMessage.getMatchLabelsMap().set(key, value)
    );
  }

  return selectorMessage;
}

function dataSourceMessageFromObject(
  dataSourceObj: DataSource.AsObject,
  dataSourceMessage = new DataSource()
): DataSource {
  const { fetchUrl, inlineBytes, inlineString, configMap } = dataSourceObj;
  if (fetchUrl !== undefined) {
    dataSourceMessage.setFetchUrl(fetchUrl);
  }
  if (inlineBytes !== undefined) {
    if (typeof inlineBytes === 'string') {
      dataSourceMessage.setInlineBytes(inlineBytes);
    } else {
      let inlineBytesUint8Array = new Uint8Array(inlineBytes);
      dataSourceMessage.setInlineBytes(inlineBytesUint8Array);
    }
  }
  if (inlineString !== undefined) {
    dataSourceMessage.setInlineString(inlineString);
  }

  if (configMap !== undefined) {
    const { name, namespace, key } = configMap;
    let newConfigMap = new DataSource.ConfigMapData();
    newConfigMap.setName(name);
    newConfigMap.setNamespace(namespace);
    newConfigMap.setKey(key);
    dataSourceMessage.setConfigMap(newConfigMap);
  }

  return dataSourceMessage;
}

function createPortal(
  portalWriteRequest: PortalWriteRequest.AsObject
): Promise<Portal.AsObject> {
  const { portal, usersList, apiDocsList, groupsList } = portalWriteRequest;
  let request = new PortalWriteRequest();
  if (portal !== undefined) {
    let portalToCreate = portalMessageFromObject(portalWriteRequest.portal!);
    request.setPortal(portalToCreate);
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
        console.log('statusMessage', statusMessage);
        if (status !== grpc.Code.OK) {
          reject(statusMessage);
        }
      }
    });
  });
}

function updatePortal(
  portalWriteRequest: PortalWriteRequest.AsObject
): Promise<Portal.AsObject> {
  const { portal, usersList, apiDocsList, groupsList } = portalWriteRequest;
  let request = new PortalWriteRequest();
  return new Promise((resolve, reject) => {
    if (portal !== undefined && portal.metadata) {
      let portalRef = new ObjectRef();
      portalRef.setName(portal.metadata.name);
      portalRef.setNamespace(portal.metadata.namespace);

      grpc.unary(PortalApi.GetPortalWithAssets, {
        request: portalRef,
        host,
        onEnd: endMessage => {
          let portalToUpdate = endMessage.message as Portal;
          let portalSpec = portalToUpdate.getSpec();

          if (
            portal.spec?.domainsList !== [] &&
            portal.spec?.domainsList !== undefined
          ) {
            portalSpec?.setDomainsList(portal.spec?.domainsList);
          }
          let portalStyling = portalSpec?.getCustomStyling();
          if (!!portal.spec?.customStyling?.backgroundColor) {
            portalStyling?.setBackgroundColor(
              portal.spec.customStyling.backgroundColor
            );
          }
          if (!!portal.spec?.customStyling?.defaultTextColor) {
            portalStyling?.setDefaultTextColor(
              portal.spec.customStyling.defaultTextColor
            );
          }
          if (!!portal.spec?.customStyling?.primaryColor) {
            portalStyling?.setPrimaryColor(
              portal.spec.customStyling.primaryColor
            );
          }
          if (!!portal.spec?.customStyling?.secondaryColor) {
            portalStyling?.setSecondaryColor(
              portal.spec.customStyling.secondaryColor
            );
          }
          if (portal.spec?.banner?.inlineBytes !== undefined) {
            let bannerMessage = portalSpec?.getBanner();
            bannerMessage?.setInlineBytes(portal.spec.banner.inlineBytes);
            portalSpec?.setBanner(bannerMessage);
          }
          if (portal.spec?.favicon?.inlineBytes !== undefined) {
            let faviconMessage = portalSpec?.getFavicon();
            faviconMessage?.setInlineBytes(portal.spec?.favicon?.inlineBytes);
            portalSpec?.setFavicon(faviconMessage);
          }
          if (portal.spec?.primaryLogo?.inlineBytes !== undefined) {
            let primaryLogoMessage = portalSpec?.getPrimaryLogo();
            primaryLogoMessage?.setInlineBytes(
              portal.spec.primaryLogo.inlineBytes
            );
            portalSpec?.setPrimaryLogo(primaryLogoMessage);
          }
          if (!!portal.spec?.description) {
            portalSpec?.setDescription(portal.spec.description);
          }
          if (!!portal.spec?.displayName) {
            portalSpec?.setDisplayName(portal.spec.displayName);
          }
          portalSpec?.setCustomStyling(portalStyling);
          portalToUpdate.setSpec(portalSpec);

          request.setPortal(portalToUpdate);

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

          grpc.invoke(PortalApi.UpdatePortal, {
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
                console.log('statusMessage', statusMessage);
                reject(statusMessage);
              }
            }
          });
        }
      });
    }
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

function createPortalPage(
  portalRef: ObjectRef.AsObject,
  staticPage: StaticPage.AsObject
): Promise<Portal.AsObject> {
  return new Promise((resolve, reject) => {
    const requestObjectRef = new ObjectRef();
    requestObjectRef.setName(portalRef.name);
    requestObjectRef.setNamespace(portalRef.namespace);

    grpc.unary(PortalApi.GetPortalWithAssets, {
      request: requestObjectRef,
      host,
      metadata: new grpc.Metadata(),
      onEnd: endMessage => {
        let portal = endMessage.message as Portal;

        let request = new PortalWriteRequest();
        let portalSpec = portal.getSpec();
        console.log(portalSpec?.toObject());

        let staticPageClass = new StaticPage();
        staticPageClass.setName(staticPage.name);
        staticPageClass.setDescription(staticPage.description);
        staticPageClass.setNavigationLinkName(staticPage.navigationLinkName);
        staticPageClass.setPath(staticPage.path);
        staticPageClass.setContent(
          createDataSourceClassFromObject(staticPage.content)
        );
        staticPageClass.setDisplayOnHomepage(staticPage.displayOnHomepage);

        portalSpec?.setStaticPagesList([
          ...portalSpec.getStaticPagesList(),
          staticPageClass
        ]);

        portal.setSpec(portalSpec);

        request.setPortal(portal);
        request.setPortalOnly(true);

        console.log(request.getPortal()?.toObject());

        grpc.invoke(PortalApi.UpdatePortal, {
          request: request,
          host,
          metadata: new grpc.Metadata(),
          onHeaders: (headers: grpc.Metadata) => {
            // console.log('onheaders', headers);
          },
          onMessage: (message: Portal) => {
            // console.log('message', message);
            if (message) {
              resolve(message.toObject());
            }
          },
          onEnd: (
            status: grpc.Code,
            statusMessage: string,
            trailers: grpc.Metadata
          ) => {
            // console.log('onEnd', status, statusMessage, trailers);
            if (status !== grpc.Code.OK) {
              reject(statusMessage);
            }
          }
        });
      }
    });
  });
}

function updatePortalPage(
  portalRef: ObjectRef.AsObject,
  staticPage: StaticPage.AsObject
): Promise<Portal.AsObject> {
  return new Promise((resolve, reject) => {
    const requestObjectRef = new ObjectRef();
    requestObjectRef.setName(portalRef.name);
    requestObjectRef.setNamespace(portalRef.namespace);

    grpc.unary(PortalApi.GetPortalWithAssets, {
      request: requestObjectRef,
      host,
      metadata: new grpc.Metadata(),
      onEnd: endMessage => {
        let portal = endMessage.message as Portal;

        let staticPageClass = new StaticPage();
        staticPageClass.setName(staticPage.name);
        staticPageClass.setDescription(staticPage.description);
        staticPageClass.setNavigationLinkName(staticPage.navigationLinkName);
        staticPageClass.setPath(staticPage.path);
        staticPageClass.setContent(
          createDataSourceClassFromObject(staticPage.content)
        );
        staticPageClass.setDisplayOnHomepage(staticPage.displayOnHomepage);

        let request = new PortalWriteRequest();
        let portalSpec = portal.getSpec();

        portalSpec?.setStaticPagesList([
          ...portalSpec
            .getStaticPagesList()
            .filter(spc => spc.getName() !== staticPageClass.getName()),
          staticPageClass
        ]);

        portal.setSpec(portalSpec);

        request.setPortal(portal);
        request.setPortalOnly(true);

        grpc.invoke(PortalApi.UpdatePortal, {
          request: request,
          host,
          metadata: new grpc.Metadata(),
          onHeaders: (headers: grpc.Metadata) => {
            // console.log('onheaders', headers);
          },
          onMessage: (message: Portal) => {
            // console.log('message', message);
            if (message) {
              resolve(message.toObject());
            }
          },
          onEnd: (
            status: grpc.Code,
            statusMessage: string,
            trailers: grpc.Metadata
          ) => {
            // console.log('onEnd', status, statusMessage, trailers);
            if (status !== grpc.Code.OK) {
              reject(statusMessage);
            }
          }
        });
      }
    });
  });
}

function deletePortalPage(
  portalRef: ObjectRef.AsObject,
  deletingPortalName: string
): Promise<Portal.AsObject> {
  return new Promise((resolve, reject) => {
    const requestObjectRef = new ObjectRef();
    requestObjectRef.setName(portalRef.name);
    requestObjectRef.setNamespace(portalRef.namespace);

    grpc.unary(PortalApi.GetPortalWithAssets, {
      request: requestObjectRef,
      host,
      metadata: new grpc.Metadata(),
      onEnd: endMessage => {
        let portal = endMessage.message as Portal;

        let request = new PortalWriteRequest();
        let portalSpec = portal.getSpec();

        portalSpec?.setStaticPagesList(
          portalSpec
            .getStaticPagesList()
            .filter(
              staticPageClass =>
                staticPageClass.getName() !== deletingPortalName
            )
        );

        portal.setSpec(portalSpec);

        request.setPortal(portal);
        request.setPortalOnly(true);

        grpc.invoke(PortalApi.UpdatePortal, {
          request: request,
          host,
          metadata: new grpc.Metadata(),
          onHeaders: (headers: grpc.Metadata) => {
            // console.log('onheaders', headers);
          },
          onMessage: (message: Portal) => {
            // console.log('message', message);
            if (message) {
              resolve(message.toObject());
            }
          },
          onEnd: (
            status: grpc.Code,
            statusMessage: string,
            trailers: grpc.Metadata
          ) => {
            // console.log('onEnd', status, statusMessage, trailers);
            if (status !== grpc.Code.OK) {
              reject(statusMessage);
            }
          }
        });
      }
    });
  });
}

function listKeyScopes(): Promise<ApiKeyScopeWithApiDocs.AsObject[]> {
  return new Promise((resolve, reject) => {
    grpc.invoke(ApiKeyScopeApi.ListApiKeyScopes, {
      request: new Empty(),
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: ApiKeyScopeList) => {
        if (message) {
          console.log(message.toObject());
          resolve(message.toObject().apiKeyScopesList);
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

function createKeyScope(
  input: ApiKeyScopeWriteRequest.AsObject
): Promise<ApiKeyScope.AsObject> {
  if (
    !input.apiKeyScope?.portal?.namespace ||
    !input.apiKeyScope?.portal?.namespace ||
    !input.apiKeyScope.spec?.namespace ||
    !input.apiKeyScope.spec?.name ||
    !input.apiKeyScope.spec?.displayName
  ) {
    return new Promise((resolve, reject) => reject('invalid request'));
  }

  let request = new ApiKeyScopeWriteRequest();
  let keyScope = new ApiKeyScope();

  let portalRef = new ObjectRef();
  portalRef.setName(input.apiKeyScope?.portal?.name);
  portalRef.setNamespace(input.apiKeyScope?.portal?.namespace);
  keyScope.setPortal(portalRef);

  let spec = new KeyScope();
  spec.setDisplayName(input.apiKeyScope.spec!.displayName);
  spec.setDescription(input.apiKeyScope.spec!.description || '');
  keyScope.setSpec(spec);
  request.setApiKeyScope(keyScope);

  let docs = input.apiDocsList.map(doc => {
    const docRef = new ObjectRef();
    docRef.setName(doc.name);
    docRef.setNamespace(doc.namespace);
    return docRef;
  });
  request.setApiDocsList(docs);

  return new Promise((resolve, reject) => {
    grpc.invoke(ApiKeyScopeApi.CreateApiKeyScope, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: ApiKeyScope) => {
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

function updateKeyScope(
  input: ApiKeyScopeWriteRequest.AsObject
): Promise<ApiKeyScope.AsObject> {
  if (
    !input.apiKeyScope?.portal?.namespace ||
    !input.apiKeyScope?.portal?.namespace ||
    !input.apiKeyScope.spec?.namespace ||
    !input.apiKeyScope.spec?.name ||
    !input.apiKeyScope.spec?.displayName
  ) {
    return new Promise((resolve, reject) => reject('invalid request'));
  }

  let request = new ApiKeyScopeWriteRequest();
  let keyScope = new ApiKeyScope();

  let portalRef = new ObjectRef();
  portalRef.setName(input.apiKeyScope?.portal?.name);
  portalRef.setNamespace(input.apiKeyScope?.portal?.namespace);
  keyScope.setPortal(portalRef);

  let spec = new KeyScope();
  spec.setName(input.apiKeyScope.spec!.name);
  spec.setDisplayName(input.apiKeyScope.spec!.displayName);
  spec.setDescription(input.apiKeyScope.spec!.description || '');
  keyScope.setSpec(spec);
  request.setApiKeyScope(keyScope);

  let docs = input.apiDocsList.map(doc => {
    const docRef = new ObjectRef();
    docRef.setName(doc.name);
    docRef.setNamespace(doc.namespace);
    return docRef;
  });
  request.setApiDocsList(docs);

  return new Promise((resolve, reject) => {
    grpc.invoke(ApiKeyScopeApi.UpdateApiKeyScope, {
      request,
      host,
      metadata: new grpc.Metadata(),
      onHeaders: (headers: grpc.Metadata) => {},
      onMessage: (message: ApiKeyScope) => {
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

function deleteKeyScope(ref: ApiKeyScopeRef.AsObject): Promise<Empty.AsObject> {
  if (!ref.portal || !ref.apiKeyScope) {
    return new Promise((resolve, reject) => reject('invalid request'));
  }

  let request = new ApiKeyScopeRef();
  let portalRef = new ObjectRef();
  portalRef.setName(ref.portal.name);
  portalRef.setNamespace(ref.portal.namespace);

  let keyScopeRef = new ObjectRef();
  keyScopeRef.setName(ref.apiKeyScope.name);
  keyScopeRef.setNamespace(ref.apiKeyScope.namespace);

  request.setPortal(portalRef);
  request.setApiKeyScope(keyScopeRef);

  return new Promise((resolve, reject) => {
    grpc.invoke(ApiKeyScopeApi.DeleteApiKeyScope, {
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
