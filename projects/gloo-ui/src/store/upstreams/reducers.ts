import { UpstreamDetails } from 'proto/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { UpstreamAction, UpstreamActionTypes } from './types';

export interface UpstreamState {
  upstreamsList: UpstreamDetails.AsObject[];
}

const initialState: UpstreamState = {
  upstreamsList: []
};

export function upstreamsReducer(
  state = initialState,
  action: UpstreamActionTypes
): UpstreamState {
  switch (action.type) {
    case UpstreamAction.LIST_UPSTREAMS:
      if (action.payload.length === state.upstreamsList.length) {
        return state;
      } else {
        return { upstreamsList: action.payload };
      }

    case UpstreamAction.UPDATE_UPSTREAM:
      return {
        ...state,
        upstreamsList: state.upstreamsList.map(upstream =>
          upstream.upstream!.metadata!.name !==
          action.payload.upstream!.metadata!.name
            ? upstream
            : action.payload
        )
      };
    case UpstreamAction.CREATE_UPSTREAM:
      return {
        ...state,
        upstreamsList: [...state.upstreamsList, action.payload]
      };
    case UpstreamAction.DELETE_UPSTREAM:
      return {
        ...state,
        upstreamsList: state.upstreamsList.filter(
          upstream =>
            upstream.upstream!.metadata!.name !== action.payload.ref!.name
        )
      };
    default:
      return state;
  }
}
