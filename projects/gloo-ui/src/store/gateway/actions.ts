import { SoloWarning } from 'Components/Common/SoloWarningContent';
import {
  UpdateGatewayRequest,
  UpdateGatewayYamlRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import { Dispatch } from 'redux';
import { MessageAction, SuccessMessageAction } from 'store/modal/types';
import { gateways } from './api';
import {
  GatewayAction,
  ListGatewaysAction,
  UpdateGatewayAction,
  UpdateGatewayYamlAction,
  UpdateGatewayYamlErrorAction
} from './types';

export const listGateways = () => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());

    try {
      const response = await gateways.listGateways();
      dispatch<ListGatewaysAction>({
        type: GatewayAction.LIST_GATEWAYS,
        payload: response.gatewayDetailsList
      });
      // dispatch(hideLoading());
    } catch (error) {}
  };
};

export const updateGateway = (
  updateGatewayRequest: UpdateGatewayRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());
    try {
      const response = await gateways.updateGateway(updateGatewayRequest);
      dispatch<UpdateGatewayAction>({
        type: GatewayAction.UPDATE_GATEWAY,
        payload: response.gatewayDetails!
      });

      dispatch<SuccessMessageAction>({
        type: MessageAction.SUCCESS_MESSAGE,
        message: 'Gateway successfully updated.'
      });
      // dispatch(hideLoading());
    } catch (error) {
      SoloWarning(
        'There was an error updating the gateway configuration.',
        error
      );
    }
  };
};

export const updateGatewayYaml = (
  updateGatewayYamlRequest: UpdateGatewayYamlRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());
    try {
      const response = await gateways.getUpdateGatewayYaml(
        updateGatewayYamlRequest
      );
      dispatch<UpdateGatewayYamlAction>({
        type: GatewayAction.UPDATE_GATEWAY_YAML,
        payload: response.gatewayDetails!
      });
    } catch (error) {
      dispatch<UpdateGatewayYamlErrorAction>({
        type: GatewayAction.UPDATE_GATEWAY_YAML_ERROR,
        payload: error
      });
      // SoloWarning(
      //   'There was an error updating the gateway configuration.',
      //   error
      // );
    }
  };
};
