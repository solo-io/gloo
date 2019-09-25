import { SoloWarning } from 'Components/Common/SoloWarningContent';
import {
  CreateUpstreamRequest,
  DeleteUpstreamRequest,
  GetUpstreamRequest,
  UpdateUpstreamRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { Dispatch } from 'redux';
import { guardByLicense } from 'store/config/actions';
import { MessageAction, SuccessMessageAction } from 'store/modal/types';
import { upstreams } from './api';
import {
  CreateUpstreamAction,
  DeleteUpstreamAction,
  GetUpstreamAction,
  ListUpstreamsAction,
  UpstreamAction
} from './types';

export const listUpstreams = () => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());
    try {
      const response = await upstreams.getUpstreamsList();
      dispatch<ListUpstreamsAction>({
        type: UpstreamAction.LIST_UPSTREAMS,
        payload: response.upstreamDetailsList
      });
      // dispatch(hideLoading());
    } catch (error) {
      // handle error
    }
  };
};

export const deleteUpstream = (
  deleteUpstreamRequest: DeleteUpstreamRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());

    try {
      guardByLicense();
      const response = await upstreams.deleteUpstream(deleteUpstreamRequest);
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
    // dispatch(showLoading());

    try {
      guardByLicense();
      const response = await upstreams.getCreateUpstream(createUpstreamRequest);
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

export const getUpstream = (
  getUpstreamRequest: GetUpstreamRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());

    try {
      const response = await upstreams.getUpstream(getUpstreamRequest);
      dispatch<GetUpstreamAction>({
        type: UpstreamAction.GET_UPSTREAM,
        payload: response.upstreamDetails!
      });
      // dispatch(hideLoading());
    } catch (error) {
      SoloWarning('There was an error retrieving the upstream.', error);
    }
  };
};

// TODO
export const updateUpstream = (
  updateUpstreamRequest: UpdateUpstreamRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    // dispatch(showLoading());
    try {
    } catch (error) {}
  };
};
