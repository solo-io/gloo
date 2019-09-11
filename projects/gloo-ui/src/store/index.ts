import { createStore, combineReducers, applyMiddleware } from 'redux';
import thunkMiddleware from 'redux-thunk';
import { composeWithDevTools } from 'redux-devtools-extension';
import { upstreamsReducer } from './upstreams/reducers';
import { virtualServicesReducer } from './virtualServices/reducers';
import { loadingBarReducer } from 'react-redux-loading-bar';
import { secretsReducer } from './secrets/reducers';
import { normalizrMiddleware } from './requests';
import { configReducer } from './config/reducers';
import { envoyReducer } from './envoy/reducers';
import { gatewaysReducer } from './gateway/reducers';
import { proxyReducer } from './proxy/reducers';
import { modalReducer } from './modal/reducers';

export const host = `${
  process.env.NODE_ENV === 'production'
    ? window.location.origin
    : 'http://localhost:8080'
  }`;

const rootReducer = combineReducers({
  upstreams: upstreamsReducer,
  virtualServices: virtualServicesReducer,
  secrets: secretsReducer,
  config: configReducer,
  gateways: gatewaysReducer,
  proxies: proxyReducer,
  envoy: envoyReducer,
  loadingBar: loadingBarReducer,
  modal: modalReducer
});

export type AppState = ReturnType<typeof rootReducer>;

export function configureStore() {
  const middlewares = [thunkMiddleware /*normalizrMiddleware */];
  const middleWareEnhancer = applyMiddleware(...middlewares);

  const store = createStore(
    rootReducer,
    composeWithDevTools(middleWareEnhancer)
  );

  return store;
}

export const globalStore = configureStore()