import * as React from 'react';
import * as jspb from 'google-protobuf';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb_service';
import { UpstreamInput } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/upstream_pb';
import { VirtualServiceInput } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/virtualservice_pb';
import { isEqual } from 'lodash';
/* -------------------------------------------------------------------------- */
/*                                   Reducer                                  */
/* -------------------------------------------------------------------------- */

// TODO: this should accept all types of variables
type VarType =
  | UpstreamInput.AsObject
  | VirtualServiceInput.AsObject
  | string[]
  | null
  | void;

type State<T> = {
  variables: VarType | null;
  prevVariables: VarType | null;
  data: T | null;
  isLoading: boolean;
  error: ServiceError | null;
  requestId: number;
};

type Action<T> =
  | { type: 'INITIAL_FETCH'; payload: null }
  | { type: 'SUCCESSFUL_FETCH'; payload: T }
  | { type: 'POLL' }
  | { type: 'FAILED_FETCH'; payload: ServiceError }
  | { type: 'VARIABLES_CHANGED'; payload: VarType }
  | { type: 'UPDATE_VARIABLES'; payload: VarType }
  | { type: 'UPDATE_VARIABLES_MANUALLY'; payload: VarType };

interface ReducerV2<T extends jspb.Message>
  extends React.Reducer<State<T>, Action<T>> {}
function reducerV2<T extends jspb.Message>(
  state: State<T>,
  action: Action<T>
): State<T> {
  switch (action.type) {
    // is this needed?
    // case 'INITIAL_FETCH':
    //   return {
    //     ...state,
    //     data: state.data,
    //     isLoading: true,
    //     error: null
    //   };
    case 'SUCCESSFUL_FETCH':
      return {
        ...state,
        variables: state.variables,
        prevVariables: state.prevVariables,
        data: action.payload,
        isLoading: false,
        error: null,
        requestId: state.requestId
      };
    case 'FAILED_FETCH':
      return {
        ...state,
        variables: state.variables,
        prevVariables: state.prevVariables,
        data: state.data,
        isLoading: false,
        error: action.payload,
        requestId: state.requestId
      };

    case 'VARIABLES_CHANGED':
      if (action.payload === state.prevVariables) {
        return state;
      }
      return {
        ...state,
        variables: action.payload,
        prevVariables: action.payload,
        data: state.data,
        isLoading: state.isLoading,
        error: null,
        requestId: 1
      };
    case 'UPDATE_VARIABLES':
      return {
        ...state,
        variables: action.payload,
        prevVariables: state.prevVariables,
        data: state.data,
        isLoading: state.isLoading,
        error: null,
        requestId: 1
      };
    case 'POLL':
      return {
        ...state,
        variables: state.variables,
        prevVariables: state.prevVariables,
        data: state.data,
        isLoading: true,
        error: null,
        requestId: state.requestId + 1
      };
    default:
      return state;
  }
}

export function useSendRequest<
  VariablesType,
  ResponseType extends jspb.Message
>(
  variables: VariablesType,
  requestFn: (variables: VariablesType) => Promise<ResponseType>,
  pollInterval = 0
) {
  const initialValue = {
    variables: variables,
    prevVariables: variables,
    isLoading: false,
    requestId: 1,
    error: null,
    data: null
  };
  const [state, dispatch] = React.useReducer<ReducerV2<ResponseType>>(
    reducerV2,
    initialValue
  );
  const updateVariables = React.useCallback(
    (variables: VariablesType) => {
      dispatch({ type: 'VARIABLES_CHANGED', payload: variables });
    },
    [dispatch]
  );

  const dispatchInitial = React.useCallback(() => {
    dispatch({ type: 'INITIAL_FETCH', payload: null });
  }, [dispatch]);

  const dispathUpdateVariables = React.useCallback(
    (variables: VariablesType) => {
      dispatch({ type: 'UPDATE_VARIABLES', payload: variables });
    },
    [dispatch]
  );

  const dispatchFetch = React.useCallback(
    data => {
      dispatch({ type: 'SUCCESSFUL_FETCH', payload: data });
    },
    [dispatch]
  );

  const dispatchError = React.useCallback(
    (error: ServiceError) => {
      dispatch({ type: 'FAILED_FETCH', payload: error });
    },
    [dispatch]
  );

  const dispatchPoll = React.useCallback(() => {
    dispatch({
      type: 'POLL'
    });
  }, [dispatch]);

  if (!isEqual(state.prevVariables, variables)) {
    updateVariables(variables);
  }
  const [params, setParams] = React.useState(variables);
  console.log(state);
  React.useEffect(() => {
    if (pollInterval > 0 && !state.isLoading) {
      console.log(state);
      const timeoutId = setTimeout(dispatchPoll, pollInterval);
      return () => {
        clearTimeout(timeoutId);
      };
    }
  }, [state.isLoading, pollInterval, dispatchPoll, params]);

  React.useEffect(() => {
    const fetchData = async () => {
      // dispatchInitial();
      try {
        const response = await requestFn(params);
        dispatchFetch(response);
      } catch (error) {
        console.error(error);
        dispatchError(error);
      }
    };
    fetchData();
  }, [params]);

  return {
    data: state.data,
    error: state.error,
    loading: state.isLoading,
    refresh: dispatchPoll,
    setNewVariables: setParams,
    update: dispathUpdateVariables,
    requestId: state.requestId
  };
}
