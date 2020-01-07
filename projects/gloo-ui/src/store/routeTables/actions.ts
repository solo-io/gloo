import { Dispatch } from 'redux';
import { routeTables } from './api';
import {
  ListRouteTablesAction,
  RouteTableAction,
  CreateRouteTableAction,
  UpdateRouteTableAction,
  UpdateRouteTableYamlAction,
  UpdateRouteTableYamlErrorAction,
  DeleteRouteTableAction
} from './types';
import { SoloWarning } from 'Components/Common/SoloWarningContent';
import {
  CreateRouteTableRequest,
  UpdateRouteTableRequest,
  UpdateRouteTableYamlRequest,
  DeleteRouteTableRequest
} from 'proto/solo-projects/projects/grpcserver/api/v1/routetable_pb';
import { guardByLicense } from 'store/config/actions';

export const listRouteTables = () => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await routeTables.listRouteTables();
      dispatch<ListRouteTablesAction>({
        type: RouteTableAction.LIST_ROUTE_TABLES,
        payload: response.routeTableDetailsList!
      });
    } catch (error) {}
  };
};

export const createRouteTable = (
  createRouteTableRequest: CreateRouteTableRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    guardByLicense();
    try {
      const response = await routeTables.createRouteTable(
        createRouteTableRequest
      );
      dispatch<CreateRouteTableAction>({
        type: RouteTableAction.CREATE_ROUTE_TABLE,
        payload: response.routeTableDetails!
      });
    } catch (error) {
      SoloWarning('There was an error creating the route table.', error);
    }
  };
};

export const updateRouteTable = (
  updateRouteTableRequest: UpdateRouteTableRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    guardByLicense();
    try {
      const response = await routeTables.updateRouteTable(
        updateRouteTableRequest
      );
      dispatch<UpdateRouteTableAction>({
        type: RouteTableAction.UPDATE_ROUTE_TABLE,
        payload: response.routeTableDetails!
      });
    } catch (error) {
      SoloWarning('There was an error updating the route table.', error);
    }
  };
};

export const updateRouteTableYaml = (
  updateRouteTableYamlRequest: UpdateRouteTableYamlRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    guardByLicense();
    try {
      const response = await routeTables.updateRouteTableYaml(
        updateRouteTableYamlRequest
      );
      dispatch<UpdateRouteTableYamlAction>({
        type: RouteTableAction.UPDATE_ROUTE_TABLE_YAML,
        payload: response.routeTableDetails!
      });
    } catch (error) {
      dispatch<UpdateRouteTableYamlErrorAction>({
        type: RouteTableAction.UPDATE_ROUTE_TABLE_YAML_ERROR,
        payload: error
      });
      //   SoloWarning('There was an error updating the route table.', error);
    }
  };
};

export const deleteRouteTable = (
  deleteRouteTableRequest: DeleteRouteTableRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    guardByLicense();
    try {
      const response = await routeTables.deleteRouteTable(
        deleteRouteTableRequest
      );
      dispatch<DeleteRouteTableAction>({
        type: RouteTableAction.DELETE_ROUTE_TABLE,
        payload: deleteRouteTableRequest
      });
    } catch (error) {
      SoloWarning('There was an error deleting the route table.', error);
    }
  };
};
