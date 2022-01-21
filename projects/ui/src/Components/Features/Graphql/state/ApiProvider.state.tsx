import * as React from 'react';
import jsonSchema from '../schema.json';

export enum ApiActions {
    UPDATE_API_URL = 'UPDATE_API_URL',
    UPDATE_API_SCHEMA = 'UPDATE_API_SCHEMA',
    UPDATE_API_NAME = 'UPDATE_API_NAME',
}

interface ApiValues {
    url: string;
    name: string;
    schema: any;
}

interface ApiAction {
    type: ApiActions;
    payload: Partial<ApiValues>;
}

const apiReducer: React.Reducer<ApiValues, ApiAction> = (
    state: ApiValues,
    action: ApiAction
): ApiValues => {
    const {
        schema = state.schema,
        url = state.url,
        name = state.name,
    } = action.payload;
    switch (action.type) {
        case ApiActions.UPDATE_API_SCHEMA:
            return {
                ...state,
                schema,
            };
        case ApiActions.UPDATE_API_URL:
            return {
                ...state,
                url,
            };
        case ApiActions.UPDATE_API_NAME:
            return {
                ...state,
                name,
            }
        default:
            return { ...state };
    }
}

export const initDefaultApiState = (
    stateOverride: Partial<ApiValues> = {}
): ApiValues => {
    return {
        url: process.env.REACT_APP_GRAPHQL_URL || 'http://localhost:4000/graphql',
        schema: jsonSchema as any,
        name: '',
        ...stateOverride,
    };
}

export const defaultState = initDefaultApiState();

type DispatchApiValues = {
    state: Partial<ApiValues>;
    dispatch: React.Dispatch<ApiAction>;
}

type ApiContextValues = ApiValues | DispatchApiValues;

const ApiContext = React.createContext<ApiContextValues>(defaultState);

export const ApiProvider: React.FC = ({ children }) => {
    const [state, dispatch] = React.useReducer(
        apiReducer,
        defaultState,
        initDefaultApiState,
    );
    const value = { state, dispatch };
    return (
        <ApiContext.Provider value={value}>{children}</ApiContext.Provider>
    )
}

export const useApiProvider = () => {
    const context = React.useContext(ApiContext);
    if (context === undefined) {
        throw new Error('Api requires a context object!');
    }
    return context as DispatchApiValues;
}
