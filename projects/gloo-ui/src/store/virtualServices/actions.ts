import { Modal } from 'antd';
import { SoloWarning } from 'Components/Common/SoloWarningContent';
import { IngressRateLimit } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/ratelimit/ratelimit_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import {
  CreateRouteRequest,
  CreateVirtualServiceRequest,
  DeleteRouteRequest,
  DeleteVirtualServiceRequest,
  ExtAuthInput,
  ShiftRoutesRequest,
  UpdateVirtualServiceRequest,
  UpdateVirtualServiceYamlRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { Dispatch } from 'redux';
import { guardByLicense } from 'store/config/actions';
import { virtualServices } from './api';
import {
  CreateRouteAction,
  CreateVirtualServiceAction,
  DeleteRouteAction,
  DeleteVirtualServiceAction,
  ListVirtualServicesAction,
  ShiftRoutesAction,
  UpdateVirtualServiceAction,
  UpdateVirtualServiceYamlAction,
  UpdateVirtualServiceYamlErrorAction,
  VirtualServiceAction
} from './types';
const { warning } = Modal;

/* --------------------------------- ACTIONS -------------------------------- */

export const listVirtualServices = () => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await virtualServices.getListVirtualServices();
      dispatch<ListVirtualServicesAction>({
        type: VirtualServiceAction.LIST_VIRTUAL_SERVICES,
        payload: response.virtualServiceDetailsList!
      });
    } catch (error) {}
  };
};

export const createVirtualService = (
  createVirtualServiceRequest: CreateVirtualServiceRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    guardByLicense();
    try {
      const response = await virtualServices.getCreateVirtualService(
        createVirtualServiceRequest
      );
      dispatch<CreateVirtualServiceAction>({
        type: VirtualServiceAction.CREATE_VIRTUAL_SERVICE,
        payload: response.virtualServiceDetails!
      });
    } catch (error) {
      warning({
        title: 'There was an error creating the virtual service.',
        content: error.message
      });
    }
  };
};

export const updateVirtualService = (
  updateVirtualServiceRequest: UpdateVirtualServiceRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await virtualServices.getUpdateVirtualService(
        updateVirtualServiceRequest
      );
      dispatch<UpdateVirtualServiceAction>({
        type: VirtualServiceAction.UPDATE_VIRTUAL_SERVICE,
        payload: response.virtualServiceDetails!
      });
    } catch (error) {
      warning({
        title: 'There was an error updating the virtual service.',
        content: error.message
      });
    }
  };
};

export const createRoute = (
  createRouteRequest: CreateRouteRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await virtualServices.getCreateRoute(createRouteRequest);
      dispatch<CreateRouteAction>({
        type: VirtualServiceAction.CREATE_ROUTE,
        payload: response.virtualServiceDetails!
      });
    } catch (error) {
      warning({
        title: 'There was an error creating the route.',
        content: error.message
      });
    }
  };
};

export const updateDomains = (updateDomainsRequest: {
  ref: ResourceRef.AsObject;
  domains: string[];
}) => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await virtualServices.getUpdateDomains(
        updateDomainsRequest
      );
      dispatch<UpdateVirtualServiceAction>({
        type: VirtualServiceAction.UPDATE_VIRTUAL_SERVICE,
        payload: response.virtualServiceDetails!
      });
    } catch (error) {
      warning({
        title: 'There was an error updating the virtual service domains.',
        content: error.message
      });
    }
  };
};

export const deleteVirtualService = (
  deleteVirtualServiceRequest: DeleteVirtualServiceRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    try {
      guardByLicense();
      const response = await virtualServices.getDeleteVirtualService(
        deleteVirtualServiceRequest
      );
      dispatch<DeleteVirtualServiceAction>({
        type: VirtualServiceAction.DELETE_VIRTUAL_SERVICE,
        payload: deleteVirtualServiceRequest
      });
    } catch (error) {
      SoloWarning('There was an error deleting the virtual service.', error);
    }
  };
};

export const updateVirtualServiceYaml = (
  updateVirtualServiceYamlRequest: UpdateVirtualServiceYamlRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await virtualServices.getUpdateVirtualServiceYaml(
        updateVirtualServiceYamlRequest
      );
      dispatch<UpdateVirtualServiceYamlAction>({
        type: VirtualServiceAction.UPDATE_VIRTUAL_SERVICE_YAML,
        payload: response.virtualServiceDetails!
      });
    } catch (error) {
      //handle error
      dispatch<UpdateVirtualServiceYamlErrorAction>({
        type: VirtualServiceAction.UPDATE_VIRTUAL_SERVICE_YAML_ERROR,
        payload: error
      });
    }
  };
};

export const deleteRoute = (
  deleteRouteRequest: DeleteRouteRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await virtualServices.getDeleteRoute(deleteRouteRequest);
      dispatch<DeleteRouteAction>({
        type: VirtualServiceAction.DELETE_ROUTE,
        payload: response.virtualServiceDetails!
      });
    } catch (error) {
      SoloWarning('There was an error deleting the route.', error);
    }
  };
};

export const shiftRoutes = (
  shiftRoutesRequest: ShiftRoutesRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await virtualServices.getShiftRoutes(shiftRoutesRequest);
      dispatch<ShiftRoutesAction>({
        type: VirtualServiceAction.SHIFT_ROUTES,
        payload: response.virtualServiceDetails!
      });
    } catch (error) {}
  };
};

export const updateRateLimit = (updateRateLimitRequest: {
  ref: ResourceRef.AsObject;
  rateLimit: IngressRateLimit.AsObject;
}) => {
  let { rateLimit, ref } = updateRateLimitRequest;
  return async (dispatch: Dispatch) => {
    try {
      const response = await virtualServices.getUpdateRateLimit({
        ref,
        rateLimitConfig: rateLimit
      });
      dispatch<UpdateVirtualServiceAction>({
        type: VirtualServiceAction.UPDATE_VIRTUAL_SERVICE,
        payload: response.virtualServiceDetails!
      });
    } catch (error) {}
  };
};

export const updateExtAuth = (updateExtAuthRequest: {
  ref: ResourceRef.AsObject;
  extAuthConfig: ExtAuthInput.AsObject;
}) => {
  let { extAuthConfig, ref } = updateExtAuthRequest;

  return async (dispatch: Dispatch) => {
    try {
      const response = await virtualServices.getUpdateExtAuth({
        ref,
        extAuthConfig
      });
      dispatch<UpdateVirtualServiceAction>({
        type: VirtualServiceAction.UPDATE_VIRTUAL_SERVICE,
        payload: response.virtualServiceDetails!
      });
    } catch (error) {}
  };
};
