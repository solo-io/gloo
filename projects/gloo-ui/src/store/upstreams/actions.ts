import { SoloWarning } from 'Components/Common/SoloWarningContent';
import {
  CreateUpstreamRequest,
  DeleteUpstreamRequest,
  UpdateUpstreamRequest
} from 'proto/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { Dispatch } from 'redux';
import { guardByLicense } from 'store/config/actions';
import { MessageAction, SuccessMessageAction } from 'store/modal/types';
import { upstreamAPI } from './api';
import {
  CreateUpstreamAction,
  DeleteUpstreamAction,
  ListUpstreamsAction,
  UpstreamAction,
  UpdateUpstreamAction
} from './types';

export const listUpstreams = () => {
  return async (dispatch: Dispatch) => {
    try {
      const response = await upstreamAPI.listUpstreams();
      dispatch<ListUpstreamsAction>({
        type: UpstreamAction.LIST_UPSTREAMS,
        payload: response
      });
    } catch (error) {}
  };
};

export const deleteUpstream = (
  deleteUpstreamRequest: DeleteUpstreamRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());

    try {
      guardByLicense();
      const response = await upstreamAPI.deleteUpstream(deleteUpstreamRequest);
      dispatch<DeleteUpstreamAction>({
        type: UpstreamAction.DELETE_UPSTREAM,
        payload: deleteUpstreamRequest
      });
      // dispatch(hideLoading());
    } catch (error) {
      SoloWarning('There was an error deleting the upstream.', error);
    }
  };
};

export const createUpstream = (
  createUpstreamRequest: CreateUpstreamRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    guardByLicense();
    try {
      const response = await upstreamAPI.createUpstream(createUpstreamRequest);
      dispatch<CreateUpstreamAction>({
        type: UpstreamAction.CREATE_UPSTREAM,
        payload: response.upstreamDetails!
      });
      // dispatch(hideLoading());
      dispatch<SuccessMessageAction>({
        type: MessageAction.SUCCESS_MESSAGE,
        message: `Upstream ${
          response.upstreamDetails!.upstream!.metadata!.name
        } successfully created.`
      });
    } catch (error) {
      SoloWarning('There was an error creating the upstream.', error);
    }
  };
};

export const updateUpstream = (
  updateUpstreamRequest: UpdateUpstreamRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());
    guardByLicense();
    try {
      const response = await upstreamAPI.updateUpstream(updateUpstreamRequest);
      dispatch<UpdateUpstreamAction>({
        type: UpstreamAction.UPDATE_UPSTREAM,
        payload: response.upstreamDetails!
      });
    } catch (error) {
      SoloWarning('There was an error updating the upstream.', error);
    }
  };
};
