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
  data: T | null;
  dataObj: any;
  isLoading: boolean;
  error: ServiceError | null;
};

type Action<T> =
  | { type: 'INITIAL_FETCH'; payload: null }
  | { type: 'SUCCESSFUL_FETCH'; payload: T }
  | { type: 'FAILED_FETCH'; payload: ServiceError };

interface ReducerV2<T extends jspb.Message>
  extends React.Reducer<State<T>, Action<T>> {}
function reducerV2<T extends jspb.Message>(
  state: State<T>,
  action: Action<T>
): State<T> {
  switch (action.type) {
    case 'INITIAL_FETCH':
      return {
        ...state,
        data: state.data,
        isLoading: true,
        error: null
      };
    case 'SUCCESSFUL_FETCH':
      return {
        ...state,
        data: action.payload,
        dataObj: action.payload.toObject(),
        isLoading: false,
        error: null
      };
    case 'FAILED_FETCH':
      return {
        ...state,
        data: state.data,
        isLoading: false,
        error: action.payload
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
  requestFn: (variables: VariablesType) => Promise<ResponseType>
) {
  const initialValue = {
    isLoading: true,
    variables: null,
    error: null,
    data: null,
    dataObj: null
  };
  const [state, dispatch] = React.useReducer<ReducerV2<ResponseType>>(
    reducerV2,
    initialValue
  );

  const [params, setParams] = React.useState(variables);

  React.useEffect(() => {
    const fetchData = async () => {
      dispatch({ type: 'INITIAL_FETCH', payload: null });
      try {
        const response = await requestFn(params);
        dispatch({ type: 'SUCCESSFUL_FETCH', payload: response });
      } catch (error) {
        console.error(error);
        dispatch({ type: 'FAILED_FETCH', payload: error as ServiceError });
      }
    };
    fetchData();
  }, [params]);

  return {
    dataObj: state.dataObj,
    dataBlob: state.data,
    data: state.data,
    error: state.error,
    loading: state.isLoading,
    setNewVariables: setParams
  };
}
