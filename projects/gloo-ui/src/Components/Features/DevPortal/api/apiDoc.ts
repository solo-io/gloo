import { ApiDocApi } from 'proto/dev-portal/api/grpc/admin/apidoc_pb_service';

import {
  ApiDoc,
  ApiDocList,
  ApiDocFilter,
  ApiDocGetRequest
} from 'proto/dev-portal/api/grpc/admin/apidoc_pb';
import { ObjectRef } from 'proto/dev-portal/api/dev-portal/v1/common_pb';
import { grpc } from '@improbable-eng/grpc-web';
import { host } from 'store';

export const apiDocApi = {
  listApiDocs,
  getApiDoc
};

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
