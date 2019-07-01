export enum RequestAction {
  START,
  SUCCESS,
  ERROR
}

interface Action<T> {
  type: RequestAction;
  payload: T;
  error?: unknown;
}

interface State<T> {
  isLoading: boolean;
  error?: unknown;
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

    default:
      throw new Error();
  }
};
