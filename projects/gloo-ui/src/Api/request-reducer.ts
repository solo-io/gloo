import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb_service';

export enum RequestAction {
  START,
  SUCCESS,
  ERROR,
  INITIALREFETCH
}

interface Action<T> {
  type: RequestAction;
  payload: T;
  error?: ServiceError;
}

interface State<T> {
  isLoading: boolean;
  error?: ServiceError;
  data: T;
}

export interface Reducer<T> extends React.Reducer<State<T>, Action<T>> {}

export const requestReducer = <T>(state: State<T>, action: Action<T>) => {
  switch (action.type) {
    case RequestAction.START:
      return {
        ...state,
        isLoading: true
      };

    case RequestAction.SUCCESS:
      return {
        ...state,
        isLoading: false,
        error: undefined,
        data: action.payload
      };

    case RequestAction.ERROR:
      return {
        ...state,
        isLoading: false
        // error: action.error
      };

    case RequestAction.INITIALREFETCH:
      return {
        ...state,
        isLoading: false,
        error: undefined
      };

    default:
      throw new Error();
  }
};
