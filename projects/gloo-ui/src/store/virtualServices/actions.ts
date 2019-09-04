import {
  ListVirtualServicesRequest,
  ListVirtualServicesResponse,
  GetVirtualServiceRequest,
  GetVirtualServiceResponse,
  DeleteVirtualServiceRequest,
  DeleteVirtualServiceResponse,
  DeleteRouteRequest,
  DeleteRouteResponse,
  SwapRoutesRequest,
  SwapRoutesResponse,
  ShiftRoutesRequest,
  ShiftRoutesResponse,
  CreateVirtualServiceRequest,
  CreateVirtualServiceResponse,
  VirtualServiceInputV2,
  UpdateVirtualServiceRequest,
  UpdateVirtualServiceResponse,
  VirtualServiceInput,
  RepeatedRoutes,
  UpdateVirtualServiceYamlRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { client } from 'Api/v2/VirtualServiceClient';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  Route,
  DirectResponseAction,
  Matcher,
  QueryParameterMatcher,
  HeaderMatcher,
  RouteAction,
  Destination
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb';
import { StringValue } from 'google-protobuf/google/protobuf/wrappers_pb';
import { Dispatch } from 'redux';
import {
  ListVirtualServicesAction,
  VirtualServiceAction,
  VirtualServiceActionTypes,
  DeleteVirtualServiceAction,
  DeleteRouteAction,
  ShiftRoutesAction,
  UpdateVirtualServiceYamlAction
} from './types';
import { showLoading, hideLoading } from 'react-redux-loading-bar';
import { EditedResourceYaml } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types_pb';
import { getResourceRef } from 'Api/v2/helpers';
import { Modal } from 'antd';
const { warning } = Modal;

export function getListVirtualServices(
  listVirtualServicesRequest: ListVirtualServicesRequest.AsObject
): Promise<ListVirtualServicesResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new ListVirtualServicesRequest();
    request.setNamespacesList(listVirtualServicesRequest.namespacesList);
    client.listVirtualServices(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        // TODO: normalize
        resolve(data!.toObject());
      }
    });
  });
}

export function getGetVirtualService(
  getVirtualServiceRequest: GetVirtualServiceRequest.AsObject
): Promise<GetVirtualServiceResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new GetVirtualServiceRequest();
    let ref = new ResourceRef();
    ref.setName(getVirtualServiceRequest.ref!.name);
    ref.setNamespace(getVirtualServiceRequest.ref!.namespace);
    request.setRef(ref);
    client.getVirtualService(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        // TODO: normalize
        resolve(data!.toObject());
      }
    });
  });
}

export function getDeleteVirtualService(
  deleteVirtualServiceRequest: DeleteVirtualServiceRequest.AsObject
): Promise<DeleteVirtualServiceResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new DeleteVirtualServiceRequest();
    let ref = new ResourceRef();
    ref.setName(deleteVirtualServiceRequest.ref!.name);
    ref.setNamespace(deleteVirtualServiceRequest.ref!.namespace);
    request.setRef(ref);
    client.deleteVirtualService(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        // TODO: normalize
        resolve(data!.toObject());
      }
    });
  });
}

export function getUpdateVirtualServiceYaml(
  updateVirtualServiceYamlRequest: UpdateVirtualServiceYamlRequest.AsObject
): Promise<UpdateVirtualServiceResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new UpdateVirtualServiceYamlRequest();

    let editedYamlData = new EditedResourceYaml();
    editedYamlData.setRef(
      getResourceRef(
        updateVirtualServiceYamlRequest.editedYamlData!.ref!.name,
        updateVirtualServiceYamlRequest.editedYamlData!.ref!.namespace
      )
    );
    editedYamlData.setEditedYaml(
      updateVirtualServiceYamlRequest.editedYamlData!.editedYaml
    );

    request.setEditedYamlData(editedYamlData);
    client.updateVirtualServiceYaml(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        // TODO: normalize
        resolve(data!.toObject());
      }
    });
  });
}

export function getDeleteRoute(
  deleteRouteRequest: DeleteRouteRequest.AsObject
): Promise<DeleteRouteResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new DeleteRouteRequest();
    let vsRef = new ResourceRef();
    vsRef.setName(deleteRouteRequest.virtualServiceRef!.name);
    vsRef.setNamespace(deleteRouteRequest.virtualServiceRef!.namespace);
    request.setVirtualServiceRef(vsRef);
    request.setIndex(deleteRouteRequest.index);
    client.deleteRoute(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        // TODO: normalize
        resolve(data!.toObject());
      }
    });
  });
}

export function getSwapRoutes(
  swapRoutesRequest: SwapRoutesRequest.AsObject
): Promise<SwapRoutesResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new SwapRoutesRequest();
    let vsRef = new ResourceRef();
    vsRef.setName(swapRoutesRequest.virtualServiceRef!.name);
    vsRef.setNamespace(swapRoutesRequest.virtualServiceRef!.namespace);

    request.setVirtualServiceRef(vsRef);
    request.setIndex1(swapRoutesRequest.index1);
    request.setIndex2(swapRoutesRequest.index2);
    client.swapRoutes(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        // TODO: normalize
        resolve(data!.toObject());
      }
    });
  });
}

export function getShiftRoutes(
  shiftRoutesRequest: ShiftRoutesRequest.AsObject
): Promise<ShiftRoutesResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new ShiftRoutesRequest();
    let vsRef = new ResourceRef();
    vsRef.setName(shiftRoutesRequest.virtualServiceRef!.name);
    vsRef.setNamespace(shiftRoutesRequest.virtualServiceRef!.namespace);

    request.setVirtualServiceRef(vsRef);
    request.setToIndex(shiftRoutesRequest.toIndex);
    request.setFromIndex(shiftRoutesRequest.fromIndex);
    client.shiftRoutes(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        // TODO: normalize
        resolve(data!.toObject());
      }
    });
  });
}

export function getUpdateVirtualService(
  updateVirtualServiceRequest: UpdateVirtualServiceRequest.AsObject
): Promise<UpdateVirtualServiceResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new UpdateVirtualServiceRequest();
  });
}
export function getCreateVirtualService(
  createVirtualSeviceRequest: CreateVirtualServiceRequest.AsObject
): Promise<CreateVirtualServiceResponse.AsObject> {
  return new Promise((resolve, reject) => {
    const { inputV2 } = createVirtualSeviceRequest;
    let request = new CreateVirtualServiceRequest();
    let ref = new ResourceRef();
    let vsInputV2 = new VirtualServiceInputV2();
    ref.setName(inputV2!.ref!.name);
    ref.setNamespace(inputV2!.ref!.namespace);
    vsInputV2.setRef(ref);
    let stringValue = new StringValue();
    stringValue.setValue(inputV2!.displayName!.value);
    if (inputV2!.routes) {
      let repeatedRoutes = new RepeatedRoutes();
      let routesList = inputV2!.routes!.valuesList.map(r => {
        let newRoute = new Route();
        if (r.directResponseAction) {
          let dra = new DirectResponseAction();
          dra.setBody(r.directResponseAction.body);
          dra.setStatus(r.directResponseAction.status);
          newRoute.setDirectResponseAction(dra);
        }
        if (r.matcher) {
          let newMatcher = new Matcher();
          newMatcher.setPrefix(r.matcher.prefix);
          newMatcher.setExact(r.matcher.exact);
          newMatcher.setRegex(r.matcher.regex);
          newMatcher.setMethodsList(r.matcher.methodsList);
          let newQueryParamsMatcherList = r.matcher.queryParametersList.map(
            qp => {
              let newQueryParamsMatcher = new QueryParameterMatcher();
              newQueryParamsMatcher.setName(qp.name);
              newQueryParamsMatcher.setValue(qp.value);
              newQueryParamsMatcher.setRegex(qp.regex);
              return newQueryParamsMatcher;
            }
          );
          newMatcher.setQueryParametersList(newQueryParamsMatcherList);

          let newHeaderMatcherList = r.matcher.headersList.map(header => {
            let newHeaderMatcher = new HeaderMatcher();
            newHeaderMatcher.setName(header.name);
            newHeaderMatcher.setRegex(header.regex);
            newHeaderMatcher.setValue(header.value);
            return newHeaderMatcher;
          });
          newMatcher.setHeadersList(newHeaderMatcherList);
          newRoute.setMatcher(newMatcher);
        }
        if (r.redirectAction) {
          //TODO
        }
        if (r.routeAction) {
          let newRouteAction = new RouteAction();
          if (r.routeAction.single) {
            const { single } = r.routeAction;

            let newDestination = new Destination();
            let usRef = new ResourceRef();
            usRef.setName(single.upstream!.name);
            usRef.setNamespace(single.upstream!.namespace);

            newDestination.setUpstream(usRef);
            // TODO
          }
        }
        if (r.routePlugins) {
          // TODO
        }
      });
      //   repeatedRoutes.setValuesList();
    }

    request.setInputV2(vsInputV2);
  });
}

export function getCreateVirtualServiceV1(
  createVirtualSeviceRequest: CreateVirtualServiceRequest.AsObject
): Promise<CreateVirtualServiceResponse.AsObject> {
  return new Promise((resolve, reject) => {
    const { input } = createVirtualSeviceRequest;
    let request = new CreateVirtualServiceRequest();
    let vsInput = new VirtualServiceInput();

    let ref = new ResourceRef();
    ref.setName(input!.ref!.name);
    ref.setNamespace(input!.ref!.namespace);
    vsInput.setRef(ref);

    vsInput.setDisplayName(input!.displayName!);
    if (input!.domainsList) {
      vsInput.setDomainsList(createVirtualSeviceRequest.input!.domainsList);
    }
    if (input!.routesList) {
      let newRoutesList = input!.routesList.map(route => {
        let r = new Route();
        // TODO
        return r;
      });
      vsInput.setRoutesList(newRoutesList);
    }
    // if (input!.)

    request.setInput(vsInput);
  });
}
// export type VirtualServiceActionTypes =
//     | CreateVirtualServiceAction
//     | UpdateVirtualServiceAction
//     | CreateRouteAction
//     | UpdateRouteAction

export const listVirtualServices = (
  listVirtualServicesRequest: ListVirtualServicesRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());

    try {
      const response = await getListVirtualServices(listVirtualServicesRequest);
      dispatch<ListVirtualServicesAction>({
        type: VirtualServiceAction.LIST_VIRTUAL_SERVICES,
        payload: response.virtualServiceDetailsList!
      });
      dispatch(hideLoading());
    } catch (error) {}
  };
};

export const deleteVirtualService = (
  deleteVirtualServiceRequest: DeleteVirtualServiceRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());

    try {
      const response = await getDeleteVirtualService(
        deleteVirtualServiceRequest
      );
      dispatch<DeleteVirtualServiceAction>({
        type: VirtualServiceAction.DELETE_VIRTUAL_SERVICE,
        payload: deleteVirtualServiceRequest
      });
      dispatch(hideLoading());
    } catch (error) {
      warning({
        title: 'There was an error deleting the virtual service.',
        content: error.message
      });
    }
  };
};

export const updateVirtualServiceYaml = (
  updateVirtualServiceYamlRequest: UpdateVirtualServiceYamlRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    dispatch(showLoading());

    try {
      const response = await getUpdateVirtualServiceYaml(
        updateVirtualServiceYamlRequest
      );
      dispatch<UpdateVirtualServiceYamlAction>({
        type: VirtualServiceAction.UPDATE_VIRTUAL_SERVICE_YAML,
        payload: response.virtualServiceDetails!
      });
      dispatch(hideLoading());
    } catch (error) {
      //handle error
    }
  };
};

export const deleteRoute = (
  deleteRouteRequest: DeleteRouteRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await getDeleteRoute(deleteRouteRequest);
      dispatch<DeleteRouteAction>({
        type: VirtualServiceAction.DELETE_ROUTE,
        payload: response.virtualServiceDetails!
      });
    } catch (error) {
      warning({
        title: 'There was an error deleting the route.',
        content: error.message
      });
    }
  };
};

export const shiftRoutes = (
  shiftRoutesRequest: ShiftRoutesRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await getShiftRoutes(shiftRoutesRequest);
      dispatch<ShiftRoutesAction>({
        type: VirtualServiceAction.SHIFT_ROUTES,
        payload: response.virtualServiceDetails!
      });
    } catch (error) {}
  };
};
