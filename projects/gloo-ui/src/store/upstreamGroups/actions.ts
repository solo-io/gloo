import { Dispatch } from 'redux';
import { guardByLicense } from 'store/config/actions';
import { upstreamGroups } from './api';
import {
  ListUpstreamGroupsAction,
  UpstreamGroupAction,
  CreateUpstreamGroupAction,
  UpdateUpstreamGroupAction,
  DeleteUpstreamGroupAction
} from './types';
import { SoloWarning } from 'Components/Common/SoloWarningContent';
import {
  CreateUpstreamGroupRequest,
  UpdateUpstreamGroupRequest,
  DeleteUpstreamGroupRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstreamgroup_pb';

export const listUpstreamGroups = () => {
  return async (dispatch: Dispatch) => {
    guardByLicense();
    try {
      const response = await upstreamGroups.listUpstreamGroups();
      dispatch<ListUpstreamGroupsAction>({
        type: UpstreamGroupAction.LIST_UPSTREAM_GROUPS,
        payload: response.upstreamGroupDetailsList!
      });
    } catch (error) {}
  };
};

export const createUpstreamGroup = (
  createUpstreamGroupRequest: CreateUpstreamGroupRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    guardByLicense();
    try {
      const response = await upstreamGroups.createUpstreamGroup(
        createUpstreamGroupRequest
      );
      dispatch<CreateUpstreamGroupAction>({
        type: UpstreamGroupAction.CREATE_UPSTREAM_GROUP,
        payload: response.upstreamGroupDetails!
      });
    } catch (error) {
      SoloWarning('There was an error creating the upstream group', error);
    }
  };
};

export const updateUpstreamGroup = (
  updateUpstreamGroupRequest: UpdateUpstreamGroupRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    guardByLicense();
    try {
      const response = await upstreamGroups.updateUpstreamGroup(
        updateUpstreamGroupRequest
      );
      dispatch<UpdateUpstreamGroupAction>({
        type: UpstreamGroupAction.UPDATE_UPSTREAM_GROUP,
        payload: response.upstreamGroupDetails!
      });
    } catch (error) {
      SoloWarning('There was an error updating the upstream group', error);
    }
  };
};

export const deleteUpstreamGroup = (
  deleteUpstreamGroupRequest: DeleteUpstreamGroupRequest.AsObject
) => {
  return async (dispatch: Dispatch) => {
    guardByLicense();
    try {
      const response = await upstreamGroups.deleteUpstreamGroup(
        deleteUpstreamGroupRequest
      );
      dispatch<DeleteUpstreamGroupAction>({
        type: UpstreamGroupAction.DELETE_UPSTREAM_GROUP,
        payload: response
      });
    } catch (error) {
      SoloWarning('There was an error deleting the upstream group', error);
    }
  };
};
