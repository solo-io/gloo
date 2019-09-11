import { Dispatch, Middleware, MiddlewareAPI } from 'redux';
import { UpstreamActionTypes } from './upstreams/types';
import { VirtualServiceActionTypes, VirtualServiceAction } from './virtualServices/types';
import { normalize } from 'normalizr';
import { upstream } from './schemas';

enum ActionStatus {
  INITIAL = 'INITIAL',
  SUCCESS = 'SUCCESS',
  FAILED = 'FAILED '
}

interface InitialAction<T> {
  type: T;
  status: ActionStatus.INITIAL;
}

interface SuccessfulAction<T, P = any> {
  type: T;
  status: ActionStatus.SUCCESS;
  payload: P;
}

interface FailedAction<T> {
  type: T;
  status: ActionStatus.FAILED;
  payload: Error;
}

export type ActionStatusType<T, P = any> =
  | InitialAction<T>
  | SuccessfulAction<T, P>
  | FailedAction<T>;

// action creators
function initialAction<T>(type: T): InitialAction<T> {
  return {
    type,
    status: ActionStatus.INITIAL
  };
}

function successfulAction<T, P>(type: T, payload: P): SuccessfulAction<T, P> {
  console.log('type, payload', type, payload);
  return {
    type,
    status: ActionStatus.SUCCESS,
    payload
  };
}

function failedAction<T>(type: T, error: Error): FailedAction<T> {
  return {
    type,
    status: ActionStatus.FAILED,
    payload: error
  };
}
type AllActions = UpstreamActionTypes | VirtualServiceActionTypes;

export const myMiddleware: Middleware<Dispatch> = ({
  dispatch
}: MiddlewareAPI) => next => (action: AllActions) => {
  console.log('action', action);
  console.log('dispatch', dispatch);
  try {
    next(action);
  } catch (error) {
    console.log('error', error);
  }
  console.log('action.type + `LOADING`', action.type + `_LOADING`);
};

export function makeActionRequest<T, P>(
  type: T,
  action: (...args: any[]) => Promise<P>,
  ...args: any[]
) {
  return async (dispatch: Dispatch) => {
    dispatch(initialAction(type));
    try {
      const payload = await action(...args);
      dispatch(successfulAction(type, payload));
    } catch (error) {
      dispatch(failedAction(type, error));
    }
  };
}

const initialState = ActionStatus.INITIAL;
export function reduceAsyncActionStatusOf<T extends string>(type: T) {
  return (
    state: ActionStatus = initialState,
    action: ActionStatusType<T>
  ): ActionStatus => {
    if (action.type === type) {
      return action.status;
    }
    return state;
  };
}

/////////////////////
export const normalizrMiddleware: Middleware<Dispatch> = ({
  dispatch
}: MiddlewareAPI) => next => (action: AllActions) => {
  console.log('action', action);
  // if (action.type.toLocaleLowerCase().includes('upstreams')) {
  //   const normalized = normalize(action.payload, [upstream]);
  //   action = Object.assign({}, action, {
  //     payload: normalized.entities
  //   });
  //   console.log('newaction', action);
  // }
  return next(action);
};